// functions regarding day boundary splitting

package ttt

import (
	"maps"
	"slices"
	"sort"
	"time"
)

// day containers keyed by date key
type DayContainerDict map[string]*DayContainer

// container of time entries
type DayContainer struct {
    // this day as a string. this should be unique
    // 2025/01/02
    DateKey string

    Entries []*TimeEntry

    // seconds
    TotalDuration int64
}

// compute the day string of a date. before hour marks the next day. if date occurs before
// this hour, it counts as the previous day (if past midnight)
func computeDate(unixTime int64,beforeHour int) string {
    var entryDate time.Time=time.Unix(unixTime,0)

    if entryDate.Hour()<beforeHour {
        entryDate=entryDate.Add(-24*time.Hour)
    }

    return entryDate.Format("2006/01/02")
}

// group time entries into corresponding day containers
func groupTimeEntries(entries []*TimeEntry,beforeHour int) []*DayContainer {
    var dayContainers DayContainerDict=DayContainerDict{}

    var entry *TimeEntry
    for _,entry = range entries {
        var entryDate string=computeDate(entry.TimeStart,beforeHour)

        var in bool
        _,in=dayContainers[entryDate]

        if !in {
            dayContainers[entryDate]=&DayContainer{
                DateKey: entryDate,
                Entries: []*TimeEntry{},
                TotalDuration: 0,
            }
        }

        dayContainers[entryDate].Entries=append(
            dayContainers[entryDate].Entries,
            entry,
        )

        if entry.Duration>0 {
            dayContainers[entryDate].TotalDuration+=entry.Duration
        }
    }

    return slices.Collect(maps.Values(dayContainers))
}

// sort day containers in place based on date key
func sortDayContainers(dayContainers []*DayContainer) {
    sort.Slice(dayContainers, func(i int, j int) bool {
        // Parse the date keys
        var di time.Time
        var dj time.Time
        var e1 error
        var e2 error

        di, e1 = time.Parse("2006/01/02", dayContainers[i].DateKey)
        dj, e2 = time.Parse("2006/01/02", dayContainers[j].DateKey)

        if e1 != nil || e2 != nil {
            // If parsing fails, fallback to string compare descending
            return dayContainers[i].DateKey > dayContainers[j].DateKey
        }

        // Return true if i is after j (want descending order)
        return di.After(dj)
    })
}