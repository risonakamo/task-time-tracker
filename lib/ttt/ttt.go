// task time tracker package implements time tracking functionalities
package ttt

import (
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

// edit to a time entry
type TimeEntryEdit struct {
    Id string `json:"id"`

    // empty to unset
    Title string `json:"title"`
    // -1 to unset
    TimeStart int64 `json:"timeStart"`
    // -1 to unset
    TimeEnd int64 `json:"timeEnd"`
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
    var e error
    var foundI int
    foundI,e=FindTimeEntryIndex(entries,id)

    if e!=nil {
        return nil,e
    }

    return entries[foundI],nil
}

// find time entry, but return index instead
func FindTimeEntryIndex(entries []*TimeEntry,id string) (int,error) {
    var entryI int
    var entry *TimeEntry
    for entryI,entry = range entries {
        if entry.Id==id {
            return entryI,nil
        }
    }

    return 0,errors.New("failed to find entry")
}

// fixes all durations on time entries in list of entries.
// mutates the entries in the given list
func RepairTimeEntries(tasks []*TimeEntry) {
    var entry *TimeEntry
    for _,entry = range tasks {
        if entry.TimeEnd>0 {
            entry.Duration=entry.TimeEnd-entry.TimeStart
        }
    }
}

// apply edits to list of time entries. returns list, but also mutates since
// it is pointer array
func ApplyTimeEntryEdits(entries []*TimeEntry,edits []TimeEntryEdit) []*TimeEntry {
    var e error

    var edit TimeEntryEdit
    for _,edit = range edits {
        var entryIndex int
        entryIndex,e=FindTimeEntryIndex(entries,edit.Id)

        if e!=nil {
            log.Warn().Err(e).Msgf("failed to find entry to edit: %s",edit.Id)
            continue
        }

        if len(edit.Title)>0 {
            entries[entryIndex].Title=edit.Title
        }

        // if timestart/timeend are both set, but time start is over time end,
        // don't change the time.
        if edit.TimeStart>0 && edit.TimeEnd>0 && edit.TimeStart>edit.TimeEnd {

        // otherwise, set whichever one is above 0
        } else {
            if edit.TimeStart>0 {
                entries[entryIndex].TimeStart=edit.TimeStart
            }

            if edit.TimeEnd>0 {
                entries[entryIndex].TimeEnd=edit.TimeEnd
            }
        }

    }

    return entries
}