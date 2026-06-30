package main

import (
	"embed"
	"task-time-tracker/lib/ttt"
	"task-time-tracker/lib/utils"

	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
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
    var app *TTTApp=&TTTApp{
        timeEntrys: []*ttt.TimeEntry{},
        currentTask: nil,
        dayContainers: []*ttt.DayContainer{},
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