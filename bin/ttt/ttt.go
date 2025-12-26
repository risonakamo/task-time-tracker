package main

import (
	"errors"
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
    // if there is actually a current task
    CurrentTaskValid bool `json:"currentTaskValid"`
    // this might empty task if none active
    CurrentTask ttt.TimeEntry `json:"currentTask"`

    AllTasks []*ttt.TimeEntry `json:"allTasks"`

    DayContainers []*ttt.DayContainer `json:"dayContainers"`
}

func main() {
    // --- consts
	var here string = utils.GetHereDirExe()
    utils.ConfigureDefaultZeroLogger()
    var e error

    var webBuildDir string=filepath.Join(here,"../../task-time-tracker-web/build")
    var dataFile string=filepath.Join(here,"data.json")
    var beforeHour int=8


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


    // --- data load
    // timeEntrys=ttt.ExampleTimeEntries1
    initialStateLoad()
    organiseTimeEntries()


    // --- routes
    // start a task. returns the newly updated state
    // if another task was running already, ends it immediately
    app.Post("/start-task",func(c fiber.Ctx) error {
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