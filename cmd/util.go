package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type Config struct {
	ConfigDir    string
	DataJSONName string
}

type Data struct {
	Day      Day       `json:day`
	Meetings []meeting `json:"meetings"`
	Tasks    []Task    `json:"tasks"`
}

type Day struct {
	WorkingTime Schedule `json:working_time`
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

type Task struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Positions []int  `json:"positions"`
}

func getDataJsonPath() string {
	return Conf.ConfigDir + Conf.DataJSONName
}

func readData() Data {
	content, err := ioutil.ReadFile(getDataJsonPath())
	if err != nil {
		fmt.Println(err, getDataJsonPath())
		return Data{}
	}
	var data Data
	json.Unmarshal(content, &data)
	return data
}

func writeData(data Data) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(getDataJsonPath(), b, 0644)
}

func InitData(day Schedule) {
	data := Data{Day{day}, []meeting{}, []Task{}}
	writeData(data)
}

func ReadDay() Day {
	return readData().Day
}

func ReadMeetings() []meeting {
	return readData().Meetings
}

func ReadTasks() []Task {
	return readData().Tasks
}

func WriteTasks(task_arr []Task) {
	data := readData()
	data.Tasks = task_arr
	writeData(data)
}

func WriteMeetings(mtg_arr []meeting) {
	data := readData()
	data.Meetings = mtg_arr
	writeData(data)
}

func AddMeeting(m meeting) {
	mtgs := ReadMeetings()
	pos := 0
	for i, mtg := range mtgs {
		if !m.Time.Start.After(mtg.Time.Start) {
			pos = i
			break
		}
		if len(mtgs)-1 == i {
			pos = -1
		}
	}
	if len(mtgs) == 0 || pos == -1 {
		mtgs = append(mtgs, m)
	} else {
		mtgs = append(mtgs[:pos+1], mtgs[pos:]...)
		mtgs[pos] = m
	}
	WriteMeetings(mtgs)
}

func AddTasks(title string, posis []int) {
	for _, posi := range posis {
		AddTask(title, posi)
	}
}
func AddTask(title string, posi int) {
	tasks := ReadTasks()
	added := false
	id := 1
	for i, t := range tasks {
		if t.Title == title {
			if !contains(t.Positions, posi) {
				t.Positions = append(t.Positions, posi)
				tasks[i] = t
			}
			added = true
			if id < t.Id {
				id = t.Id + 1
			}
		}
	}
	if !added {
		task := Task{id, title, []int{posi}}
		tasks = append(tasks, task)
	}
	WriteTasks(tasks)
}

func contains(src []int, e int) bool {
	for _, v := range src {
		if e == v {
			return true
		}
	}
	return false
}

func BuildScheduleStruct(daysString string) Schedule {
	const layout = "2006-01-02 15:04"
	days := strings.Split(daysString, "-")
	now := time.Now()
	prefix := strings.Split(now.Format(layout), " ")[0]
	start := buildDay(days[0], prefix, layout)
	end := buildDay(days[1], prefix, layout)
	return Schedule{start, end}
}

func buildDay(day string, prefix string, layout string) time.Time {
	dayString := strings.Join([]string{prefix, day}, " ")
	dayTime, _ := time.ParseInLocation(layout, dayString, time.Now().Location())
	return dayTime
}

func CheckConflictSchedule(source Schedule, target Schedule) bool {
	return !(!target.Start.Before(source.End) || !target.End.After(source.Start))
}

func PreviousSchedule(source Schedule, target Schedule) bool {
	return !target.Start.After(source.End) && target.End.After(source.Start)
}
