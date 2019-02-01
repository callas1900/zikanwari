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

type Data struct {
	Meetings []meeting `json:"meetings"`
	Tasks    []Task    `json:"tasks"`
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

func readData() Data {
	content, err := ioutil.ReadFile(Conf.DataJSONPath)
	if err != nil {
		fmt.Println(err, Conf.DataJSONPath)
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
	ioutil.WriteFile("data.json", b, 0644)
}

func InitData() {
	data := Data{[]meeting{}, []Task{}}
	writeData(data)
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

func CheckConflictSchedule(source Schedule, target Schedule) bool {
	return !(!target.Start.Before(source.End) || !target.End.After(source.Start))
}

func PreviousSchedule(source Schedule, target Schedule) bool {
	return !target.Start.After(source.End) && target.End.After(source.Start)
}
