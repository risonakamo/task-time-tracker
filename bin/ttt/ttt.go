package main

import (
	"errors"
	"os"
	"path/filepath"
	"task-time-tracker/lib/ttt"
	"task-time-tracker/lib/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/rs/zerolog/log"
)

// request body to start new task
type StartTaskReq struct {
    Title string `json:"title"`
}

// state container to send out
type TTTState struct {
    ProjectLoaded bool `json:"projectLoaded"`

    // if there is actually a current task
    CurrentTaskValid bool `json:"currentTaskValid"`
    // this might empty task if none active
    CurrentTask ttt.TimeEntry `json:"currentTask"`

    AllTasks []*ttt.TimeEntry `json:"allTasks"`

    DayContainers []*ttt.DayContainer `json:"dayContainers"`
}

// request body to change to another project data json.
// the proj name must match a file in the data folder, including
// the json extension
type ChangeProjReq struct {
    NewProjName string
}

// new persisted simpler state
type TTTState2 struct {
    // filename of last data file. might be empty
    LastDataFile string `yaml:"lastDataFile"`
}

func main() {
    // --- consts
	var here string = utils.GetHereDirExe()
    utils.ConfigureDefaultZeroLogger()
    var e error

    var webBuildDir string=filepath.Join(here,"../../task-time-tracker-web/build")
    var dataFolder string=filepath.Join(here,"data")
    var dataConfigFile string=filepath.Join(here,"data.yml")
    var dataFile string=""

    // 24 hour time of day used to be considered the end of the day
    var beforeHour int=8

    var projectLoaded bool=false



    // --- app setup
    var app *fiber.App=fiber.New(fiber.Config{
        CaseSensitive: true,
        ErrorHandler: func(c fiber.Ctx, err error) error {
            log.Warn().Msg("fiber error")
            log.Warn().Msgf("%v",err)
            return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
        },
    })

    app.Use("/*",static.New(webBuildDir))



    // --- state
    // list of time entrys
    var timeEntrys []*ttt.TimeEntry

    // the current task. also exists in the time entrys
    var currentTask *ttt.TimeEntry=nil

    var dayContainers []*ttt.DayContainer

    var config TTTState2





    // --- functions
    // create ttt state for sending out
    createAppState:=func() TTTState {
        var currentTaskValid bool=false
        var curTaskForAppstate ttt.TimeEntry

        if currentTask!=nil {
            currentTaskValid=true
            curTaskForAppstate=*currentTask
        }

        return TTTState{
            ProjectLoaded: projectLoaded,
            CurrentTaskValid: currentTaskValid,
            CurrentTask: curTaskForAppstate,
            AllTasks: timeEntrys,
            DayContainers: dayContainers,
        }
    }

    // do operations to organise time entry related states. should be used after modifying
    // time entrys.
    // - recalculates all day containers
    // - recomputes durations on time entries
    // - sorts entrys/containers
    organiseTimeEntries:=func() {
        ttt.SortTimeEntrys(timeEntrys)
        ttt.RepairTimeEntries(timeEntrys)
        dayContainers=ttt.GroupTimeEntries(timeEntrys,beforeHour)
        ttt.SortDayContainers(dayContainers)
    }

    // read the data json and set values
    initialStateLoad:=func() {
        var savedState TTTState
        savedState,e=utils.ReadJson[TTTState](dataFile)

        if e!=nil {
            log.Warn().Err(e).Msg("failed to load data file")
            return
        }

        timeEntrys=savedState.AllTasks

        if savedState.CurrentTaskValid {
            currentTask,e=ttt.FindTimeEntry(timeEntrys,savedState.CurrentTask.Id)

            if e!=nil {
                log.Error().Msg("failed to find current task")
            }
        }
    }

    // write state to the data file
    writeState:=func() {
        e=utils.WriteJson(dataFile,createAppState())

        if e!=nil {
            log.Warn().Err(e).Msg("failed to write data file")
        }
    }

    // load the initial non project specific config file
    loadConfig:=func() TTTState2 {
        var config TTTState2
        config,e=utils.ReadYaml[TTTState2](dataConfigFile)

        if e!=nil {
            log.Warn().Err(e).Msg("failed to load config file, using empty config")
            return TTTState2{}
        }

        return config
    }

    // write to the non project specific config file
    saveConfig:=func(config TTTState2) {
        e=utils.WriteYaml(dataConfigFile,config)

        if e!=nil {
            log.Warn().Err(e).Msg("failed to save config file")
        }
    }

    // change the current project
    changeProject:=func(newProjectFile string) {
        if len(newProjectFile)==0 {
            log.Info().Msg("no last project file configured, staying unloaded")
            projectLoaded=false
            return
        }

        var targetFile string=filepath.Join(dataFolder,newProjectFile)

        _,e=os.Stat(targetFile)

        if e!=nil {
            if os.IsNotExist(e) {
                log.Warn().Msgf("last project file does not exist: %s", targetFile)
                projectLoaded=false
                return
            }
            log.Err(e).Msg("failed to stat last project file")
            projectLoaded=false
            return
        }

        dataFile=targetFile
        initialStateLoad()
        organiseTimeEntries()
        projectLoaded=true
    }





    // --- data load
    // timeEntrys=ttt.ExampleTimeEntries1
    config=loadConfig()
    changeProject(config.LastDataFile)
    if projectLoaded {
        saveConfig(config)
    }




    // --- routes
    // start a task. returns the newly updated state
    // if another task was running already, ends it immediately
    app.Post("/start-task",func(c fiber.Ctx) error {
        if !projectLoaded {
            return c.Status(fiber.StatusBadRequest).SendString("no project loaded")
        }

        var body StartTaskReq
        e=c.Bind().JSON(&body)

        if e!=nil {
            log.Err(e)
            return e
        }

        if len(body.Title)==0 {
            log.Error().Msgf("requested to start task with no title")
            return errors.New("provided no title")
        }

        if currentTask!=nil {
            ttt.EndTask(currentTask)
        }

        var newTask ttt.TimeEntry=ttt.NewTimeEntry(body.Title)

        timeEntrys=append(timeEntrys,&newTask)
        organiseTimeEntries()
        currentTask=&newTask

        var result TTTState=createAppState()
        writeState()
        return c.JSON(result)
    })

    // stops the current task, if there is any. returns the new state
    app.Post("/stop-task",func(c fiber.Ctx) error {
        if !projectLoaded {
            return c.Status(fiber.StatusBadRequest).SendString("no project loaded")
        }

        if currentTask!=nil {
            ttt.EndTask(currentTask)
            currentTask=nil
        }

        organiseTimeEntries()
        writeState()

        var result TTTState=createAppState()
        return c.JSON(result)
    })

    // get the current app state
    app.Get("/task-state",func(c fiber.Ctx) error {
        var result TTTState=createAppState()
        return c.JSON(result)
    })

    // edit a task by overwriting it. give the entire task to be edited.
    // returns the new state
    app.Post("/edit-task",func (c fiber.Ctx) error {
        if !projectLoaded {
            return c.Status(fiber.StatusBadRequest).SendString("no project loaded")
        }

        var body ttt.TimeEntry
        e=c.Bind().JSON(&body)

        if e!=nil {
            log.Err(e)
            return e
        }

        var foundEntryI int
        foundEntryI,e=ttt.FindTimeEntryIndex(timeEntrys,body.Id)

        if e!=nil {
            log.Err(e).Msg("failed to find entry to edit")
            return e
        }

        timeEntrys[foundEntryI]=&body
        organiseTimeEntries()
        writeState()

        var result TTTState=createAppState()
        return c.JSON(result)
    })

    // apply time entry edits. returns new state
    app.Post("/edit-tasks2",func (c fiber.Ctx) error {
        if !projectLoaded {
            return c.Status(fiber.StatusBadRequest).SendString("no project loaded")
        }

        var body []ttt.TimeEntryEdit
        e=c.Bind().JSON(&body)

        if e!=nil {
            log.Err(e)
            return e
        }

        timeEntrys=ttt.ApplyTimeEntryEdits(timeEntrys,body)
        organiseTimeEntries()
        writeState()

        var result TTTState=createAppState()
        return c.JSON(result)
    })

    // open the data dir
    app.Get("/open-data-folder",func(c fiber.Ctx) error {
        utils.OpenTargetWithDefaultProgram(here)

        return c.SendStatus(fiber.StatusOK)
    })

    // close the program
    app.Get("/close",func(c fiber.Ctx) error {
        return app.Shutdown()
    })

    // change to another project
    app.Post("/change-project",func(c fiber.Ctx) error {
        var body ChangeProjReq
        e=c.Bind().JSON(&body)

        if e!=nil {
            log.Err(e)
            return e
        }

        changeProject(body.NewProjName)

        if !projectLoaded {
            return c.Status(fiber.StatusBadRequest).SendString("project file does not exist")
        }

        config.LastDataFile=body.NewProjName
        saveConfig(config)

        var result TTTState=createAppState()
        return c.JSON(result)
    })

    // --- running
    e=utils.OpenTargetWithDefaultProgram(
        "http://localhost:4602",
    )

    if e!=nil {
        panic(e)
    }

    e=app.Listen(":4602")

    if e!=nil {
        panic(e)
    }
}