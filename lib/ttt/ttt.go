// task time tracker package implements time tracking functionalities
package ttt

import (
	"time"

	"github.com/google/uuid"
)

// a single task time entry
type TimeEntry struct {
    Id string
    Title string

    // unix time seconds
    TimeStart int64
    // unix time seconds. -1 if not set
    TimeEnd int64
    // seconds. -1 if not ended
    Duration int64
}

// create a new time entry with date starting at now
func NewTimeEntry(title string) TimeEntry {
    return TimeEntry{
        Id: uuid.New().String(),
        Title: title,
        TimeStart: time.Now().Unix(),
        TimeEnd: -1,
        Duration: -1,
    }
}

// given a time entry task, ends it at the current time. computes
// the duration. mutates the given entry
func EndTask(task *TimeEntry) {
    var now int64=time.Now().Unix()
    task.TimeEnd=now
    task.Duration=now-task.TimeStart
}