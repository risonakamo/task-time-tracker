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
}

func main() {
    // --- consts
	var here string = utils.GetHereDirExe()
    var e error

    var webBuildDir string=filepath.Join(here,"../../task-time-tracker-web/build")


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
    var timeEntrys []*ttt.TimeEntry=[]*ttt.TimeEntry{}

    // the current task. also exists in the time entrys
    var currentTask *ttt.TimeEntry=nil


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
        }
    }


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
        currentTask=&newTask

        var result TTTState=createAppState()
        return c.JSON(result)
    })

    // stops the current task, if there is any. returns the new state
    app.Post("/stop-task",func(c fiber.Ctx) error {
        if currentTask!=nil {
            ttt.EndTask(currentTask)
            currentTask=nil
        }

        var result TTTState=createAppState()
        return c.JSON(result)
    })

    // get the current app state
    app.Get("/task-state",func(c fiber.Ctx) error {
        var result TTTState=createAppState()
        return c.JSON(result)
    })


    // --- running
    e=app.Listen(":4602")

    if e!=nil {
        panic(e)
    }
}