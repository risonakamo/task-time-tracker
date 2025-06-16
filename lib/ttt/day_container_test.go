package ttt

import (
	"fmt"
	"testing"
	"time"

	"github.com/k0kubun/pp/v3"
)

func Test_timeEntryDate(t *testing.T) {
    var dates []time.Time=[]time.Time{
        time.Date(2025, 7, 1, 8, 0, 0, 0, time.Local),
        time.Date(2025, 7, 1, 6, 0, 0, 0, time.Local),
        time.Date(2025, 7, 1, 4, 0, 0, 0, time.Local),
        time.Date(2025, 7, 1, 10, 0, 0, 0, time.Local),
        time.Date(2025, 6, 30, 8, 0, 0, 0, time.Local),
        time.Date(2025, 6, 30, 7, 0, 0, 0, time.Local),
        time.Date(2025, 6, 29, 8, 0, 0, 0, time.Local),
        time.Date(2025, 6, 29, 6, 0, 0, 0, time.Local),
        time.Date(2025, 6, 28, 6, 0, 0, 0, time.Local),
        time.Date(2025, 6, 28, 2, 0, 0, 0, time.Local),
    }

    var beforeHour int=8

    var date time.Time
    for _,date = range dates {
        fmt.Println(
            date.String(),"->",
            computeDate(date.Unix(),beforeHour),
        )
    }
}

func Test_groupTimeEntries(t *testing.T) {
    var entries []*TimeEntry = []*TimeEntry{
        {
            Id:        "1",
            Title:     "Morning Task",
            TimeStart: time.Date(2025, 6, 10, 7, 30, 0, 0, time.Local).Unix(), // Before 8 AM
            TimeEnd:   time.Date(2025, 6, 10, 8, 30, 0, 0, time.Local).Unix(),
            Duration:  3600,
        },
        {
            Id:        "2",
            Title:     "Afternoon Task",
            TimeStart: time.Date(2025, 6, 10, 14, 0, 0, 0, time.Local).Unix(),
            TimeEnd:   time.Date(2025, 6, 10, 16, 0, 0, 0, time.Local).Unix(),
            Duration:  7200,
        },
        {
            Id:        "3",
            Title:     "Overnight Task",
            TimeStart: time.Date(2025, 6, 9, 23, 0, 0, 0, time.Local).Unix(),
            TimeEnd:   time.Date(2025, 6, 10, 1, 0, 0, 0, time.Local).Unix(),
            Duration:  7200,
        },
        {
            Id:        "4",
            Title:     "Unfinished Task",
            TimeStart: time.Date(2025, 6, 11, 10, 0, 0, 0, time.Local).Unix(),
            TimeEnd:   -1,
            Duration:  -1,
        },
        {
            Id:        "5",
            Title:     "Another Morning Task",
            TimeStart: time.Date(2025, 6, 12, 6, 45, 0, 0, time.Local).Unix(),
            TimeEnd:   time.Date(2025, 6, 12, 7, 15, 0, 0, time.Local).Unix(),
            Duration:  1800,
        },
    }

    var beforeHour int=8

    result:=GroupTimeEntries(entries,beforeHour)
    SortDayContainers(result)

    pp.Println(result)
}