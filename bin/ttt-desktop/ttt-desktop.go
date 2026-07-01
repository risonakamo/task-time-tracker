package main

import (
	"context"
	"embed"
	"errors"
	"path/filepath"
	"task-time-tracker/lib/ttt"
	"task-time-tracker/lib/utils"

	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:web-build
var MyAssets embed.FS
var BeforeHour int=8

func main() {
	var e error

    var app *TTTApp=newTTTApp()

	e=wails.Run(&options.App{
        Title: "Task Time Tracker Desktop",
		Width:  1024,
		Height: 768,
        AssetServer: &assetserver.Options{
            Assets: MyAssets,
        },
        Bind: []interface{}{
            app,
        },
    })

    if e!=nil {
        panic(e)
    }
}



type TTTApp struct {
    // list of time entrys
    timeEntrys []*ttt.TimeEntry
    // the current task. also exists in the time entrys
    currentTask *ttt.TimeEntry
    dayContainers []*ttt.DayContainer

    dataFile string

    ctx context.Context
}

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

// ttt app constructor
func newTTTApp() *TTTApp {
    var here string = utils.GetHereDirExe()

    var app *TTTApp=&TTTApp{
        timeEntrys: []*ttt.TimeEntry{},
        currentTask: nil,
        dayContainers: []*ttt.DayContainer{},
        dataFile: filepath.Join(here, "data.json"),
    }

    app.initialStateLoad()
    app.organiseTimeEntries()

    return app
}

// create current ttt state for sending out
func (app *TTTApp) createAppState() TTTState {
    var currentTaskValid bool=false
    var curTaskForAppstate ttt.TimeEntry

    if app.currentTask!=nil {
        currentTaskValid=true
        curTaskForAppstate=*app.currentTask
    }

    return TTTState{
        CurrentTaskValid: currentTaskValid,
        CurrentTask: curTaskForAppstate,
        AllTasks: app.timeEntrys,
        DayContainers: app.dayContainers,
    }
}

// do operations to organise time entry related states. should be used after modifying
// time entrys.
// - recalculates all day containers
// - recomputes durations on time entries
// - sorts entrys/containers
func (app *TTTApp) organiseTimeEntries() {
    ttt.SortTimeEntrys(app.timeEntrys)
    ttt.RepairTimeEntries(app.timeEntrys)
    app.dayContainers=ttt.GroupTimeEntries(app.timeEntrys,BeforeHour)
    ttt.SortDayContainers(app.dayContainers)
}

// read the data json and set values
func (app *TTTApp) initialStateLoad() {
    var e error
    var savedState TTTState
    savedState,e=utils.ReadJson[TTTState](app.dataFile)

    if e!=nil {
        log.Warn().Err(e).Msg("failed to load data file")
        return
    }

    app.timeEntrys=savedState.AllTasks

    if savedState.CurrentTaskValid {
        app.currentTask,e=ttt.FindTimeEntry(app.timeEntrys,savedState.CurrentTask.Id)

        if e!=nil {
            log.Error().Msg("failed to find current task")
        }
    }
}

// write state to the data file
func (app *TTTApp) writeState() {
    var e error
    e=utils.WriteJson(app.dataFile,app.createAppState())

    if e!=nil {
        log.Warn().Err(e).Msg("failed to write data file")
    }
}

// start a task. returns the newly updated state
// if another task was running already, ends it immediately
func (app *TTTApp) StartTask(startTaskReq StartTaskReq) (TTTState,error) {
    if len(startTaskReq.Title)==0 {
        log.Error().Msgf("requested to start task with no title")
        return TTTState{},errors.New("provided no title")
    }

    if app.currentTask!=nil {
        ttt.EndTask(app.currentTask)
    }

    var newTask ttt.TimeEntry=ttt.NewTimeEntry(startTaskReq.Title)

    app.timeEntrys=append(app.timeEntrys,&newTask)
    app.organiseTimeEntries()
    app.currentTask=&newTask

    var result TTTState=app.createAppState()
    app.writeState()
    return result,nil
}

// stop the current task, if there is any. returns the new state
func (app *TTTApp) StopTask() (TTTState,error) {
    if app.currentTask!=nil {
        ttt.EndTask(app.currentTask)
        app.currentTask=nil
    }

    app.organiseTimeEntries()
    app.writeState()

    var result TTTState=app.createAppState()
    return result,nil
}

// get the current app state
func (app *TTTApp) TaskState() TTTState {
    return app.createAppState()
}

// edit a task by overwriting it. give the entire task to be edited.
// returns the new state
func (app *TTTApp) EditTask(body ttt.TimeEntry) (TTTState,error) {
    var foundEntryI int
    var e error
    foundEntryI,e=ttt.FindTimeEntryIndex(app.timeEntrys,body.Id)

    if e!=nil {
        log.Err(e).Msg("failed to find entry to edit")
        return TTTState{},e
    }

    app.timeEntrys[foundEntryI]=&body
    app.organiseTimeEntries()
    app.writeState()

    var result TTTState=app.createAppState()
    return result,nil
}

// apply time entry edits. returns new state
func (app *TTTApp) EditTasks2(body []ttt.TimeEntryEdit) (TTTState,error) {
    app.timeEntrys=ttt.ApplyTimeEntryEdits(app.timeEntrys,body)
    app.organiseTimeEntries()
    app.writeState()

    var result TTTState=app.createAppState()
    return result,nil
}

// open the data dir
func (app *TTTApp) OpenDataFolder() error {
    return utils.OpenTargetWithDefaultProgram(filepath.Dir(app.dataFile))
}

// close the program
func (app *TTTApp) Close() {
    runtime.Quit(app.ctx)
}

// wails startup func
func (app *TTTApp) startup(ctx context.Context) {
    app.ctx=ctx
}