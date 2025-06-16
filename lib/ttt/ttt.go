// task time tracker package implements time tracking functionalities
package ttt

import (
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
)

// a single task time entry
type TimeEntry struct {
    Id string `json:"id"`
    Title string `json:"title"`

    // unix time seconds
    TimeStart int64 `json:"timeStart"`
    // unix time seconds. -1 if not set
    TimeEnd int64 `json:"timeEnd"`
    // seconds. -1 if not ended
    Duration int64 `json:"duration"`
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

// sorts time entry list by start time (latest comes first)
// mutates the input array
func SortTimeEntrys(tasks []*TimeEntry) {
    sort.Slice(tasks, func(task1i int,task2i int) bool {
        return tasks[task1i].TimeStart > tasks[task2i].TimeStart
    })
}

// find pointer to a time entry from list
func FindTimeEntry(entries []*TimeEntry,id string) (*TimeEntry,error) {
    var entry *TimeEntry
    for _,entry = range entries {
        if entry.Id==id {
            return entry,nil
        }
    }

    return nil,errors.New("failed to find entry")
}