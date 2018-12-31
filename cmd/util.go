package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

type Schedule struct {
	Start time.Time
	End   time.Time
}

func ReadMeetings() []meeting {
	content, err := ioutil.ReadFile(Conf.DataJSONPath)
	if err != nil {
		log.Fatal(err)
	}
	var mtgs meetings
	json.Unmarshal(content, &mtgs)
	return mtgs.Meetings
}

func CheckConflictSchedule(source Schedule, target Schedule) bool {
	return !target.End.After(source.Start) || !target.Start.Before(source.End)
}

func PreviousSchedule(source Schedule, target Schedule) bool {
	return !target.Start.After(source.End) && target.End.After(source.Start)
}
