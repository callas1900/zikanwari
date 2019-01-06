package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Config struct {
	DataJSONPath string
}

type meetings struct {
	Meetings []meeting `json:"meetings"`
}

type meeting struct {
	Id    int      `json:"id"`
	Title string   `json:"title"`
	Time  Schedule `json:"time"`
}

type Schedule struct {
	Start time.Time
	End   time.Time
}

func ReadMeetings() []meeting {
	content, err := ioutil.ReadFile(Conf.DataJSONPath)
	if err != nil {
		return []meeting{}
	}
	var mtgs meetings
	json.Unmarshal(content, &mtgs)
	return mtgs.Meetings
}

func WriteMeetings(ms meetings) {
	b, err := json.Marshal(ms)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile("data.json", b, 0644)
}

func CheckConflictSchedule(source Schedule, target Schedule) bool {
	return !(!target.Start.Before(source.End) || !target.End.After(source.Start))
}

func PreviousSchedule(source Schedule, target Schedule) bool {
	return !target.Start.After(source.End) && target.End.After(source.Start)
}
