// example time entries

package ttt

import (
	"time"

	"github.com/google/uuid"
)

// some time entries
var ExampleTimeEntries1 []*TimeEntry=[]*TimeEntry{
    // June 12, 2025
    {
        Id:        uuid.New().String(),
        Title:     "Morning meeting",
        TimeStart: time.Date(2025, 6, 12, 9, 0, 0, 0, time.Local).Unix(),
        TimeEnd:   time.Date(2025, 6, 12, 10, 0, 0, 0, time.Local).Unix(),
        Duration:  3600,
    },
    {
        Id:        uuid.New().String(),
        Title:     "Emails and planning",
        TimeStart: time.Date(2025, 6, 12, 10, 15, 0, 0, time.Local).Unix(),
        TimeEnd:   time.Date(2025, 6, 12, 11, 30, 0, 0, time.Local).Unix(),
        Duration:  4500,
    },
    {
        Id:        uuid.New().String(),
        Title:     "Code review",
        TimeStart: time.Date(2025, 6, 12, 13, 0, 0, 0, time.Local).Unix(),
        TimeEnd:   time.Date(2025, 6, 12, 14, 0, 0, 0, time.Local).Unix(),
        Duration:  3600,
    },

    // June 13, 2025
    {
        Id:        uuid.New().String(),
        Title:     "Sprint planning",
        TimeStart: time.Date(2025, 6, 13, 9, 30, 0, 0, time.Local).Unix(),
        TimeEnd:   time.Date(2025, 6, 13, 10, 30, 0, 0, time.Local).Unix(),
        Duration:  3600,
    },
    {
        Id:        uuid.New().String(),
        Title:     "Feature implementation",
        TimeStart: time.Date(2025, 6, 13, 11, 0, 0, 0, time.Local).Unix(),
        TimeEnd:   time.Date(2025, 6, 13, 13, 30, 0, 0, time.Local).Unix(),
        Duration:  9000,
    },
    {
        Id:        uuid.New().String(),
        Title:     "Bug fixing",
        TimeStart: time.Date(2025, 6, 13, 14, 0, 0, 0, time.Local).Unix(),
        TimeEnd:   time.Date(2025, 6, 13, 15, 15, 0, 0, time.Local).Unix(),
        Duration:  4500,
    },

    // June 14, 2025
    {
        Id:        uuid.New().String(),
        Title:     "Testing",
        TimeStart: time.Date(2025, 6, 14, 9, 0, 0, 0, time.Local).Unix(),
        TimeEnd:   time.Date(2025, 6, 14, 10, 30, 0, 0, time.Local).Unix(),
        Duration:  5400,
    },
    {
        Id:        uuid.New().String(),
        Title:     "Documentation",
        TimeStart: time.Date(2025, 6, 14, 10, 45, 0, 0, time.Local).Unix(),
        TimeEnd:   time.Date(2025, 6, 14, 11, 30, 0, 0, time.Local).Unix(),
        Duration:  2700,
    },
    {
        Id:        uuid.New().String(),
        Title:     "Demo prep",
        TimeStart: time.Date(2025, 6, 14, 13, 0, 0, 0, time.Local).Unix(),
        TimeEnd:   time.Date(2025, 6, 14, 14, 30, 0, 0, time.Local).Unix(),
        Duration:  5400,
    },
    {
        Id:        uuid.New().String(),
        Title:     "Team sync",
        TimeStart: time.Date(2025, 6, 14, 15, 0, 0, 0, time.Local).Unix(),
        TimeEnd:   time.Date(2025, 6, 14, 15, 30, 0, 0, time.Local).Unix(),
        Duration:  1800,
    },
}