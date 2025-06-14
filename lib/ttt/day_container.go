// functions regarding day boundary splitting

package ttt

import "time"

// container of time entries
type DayContainer struct {
    Id string

    // this day as a string
    // 2025/01/02
    DateKey string

    Entries []*TimeEntry
    TotalDuration int64
}

// compute the day string of a date. before hour marks the next day. if task occurs before
// this hour, it counts as the previous day (if past midnight)
func computeDate(unixTime int64,beforeHour int) string {
    var entryDate time.Time=time.Unix(unixTime,0)

    if entryDate.Hour()<beforeHour {
        entryDate=entryDate.Add(-24*time.Hour)
    }

    return entryDate.Format("2006/01/02")
}