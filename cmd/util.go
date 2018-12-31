package cmd

import (
	"time"
)

type Schedule struct {
	Start time.Time
	End   time.Time
}

func CheckConflictSchedule(source Schedule, target Schedule) bool {
	return !target.End.After(source.Start) || !target.Start.Before(source.End)
}

func PreviousSchedule(source Schedule, target Schedule) bool {
	return !target.Start.After(source.End) && target.End.After(source.Start)
}
