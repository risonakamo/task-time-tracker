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
    CurrentTask ttt.TimeEntry `json:"currentTask"`
    AllTasks []ttt.TimeEntry `json:"allTasks"`
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
    var timeEntrys []ttt.TimeEntry=[]ttt.TimeEntry{}
    var currentTask *ttt.TimeEntry=nil


    // --- functions
    // create ttt state for sending out
    createAppState:=func() TTTState {
        return TTTState{
            CurrentTask: *currentTask,
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

        timeEntrys=append(timeEntrys,newTask)
        currentTask=&newTask

        var result TTTState=createAppState()
        return c.JSON(result)
    })


    // --- running
    e=app.Listen(":4602")

    if e!=nil {
        panic(e)
    }
}