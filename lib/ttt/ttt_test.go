package ttt

import (
	"fmt"
	"testing"
	"time"
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