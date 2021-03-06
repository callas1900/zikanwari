// Copyright © 2018 RYO TANAKA <callas1900@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const time_layout = "15:04"

type Pomo struct {
	Id   int      `json:"id"`
	Name string   `json:"title"`
	Time Schedule `json:"time"`
}

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show pomos and meetings",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		mtgs := ReadMeetings()
		// TODO: remove dummy impl
		day := ReadDay()
		pomos := CalcPomos(mtgs, day.WorkingTime.Start, day.WorkingTime.End, Conf.pomoTime, Conf.pomoRest)
		if cmd.Flag("verbose").Changed {
			displayHeader()
		}
		display(mtgs, pomos)
		if cmd.Flag("verbose").Changed {
			displayFooter()
		}
	},
}

func CalcPomos(mtgs []meeting, start time.Time, end time.Time, unit int, rest int) []Pomo {
	maxPomoCount := int(math.Round(end.Sub(start).Minutes() / float64(unit+rest)))
	pomos := make([]Pomo, maxPomoCount)
	s := start
	e := s.Add(time.Duration(unit) * time.Minute)
	count := 0
	for i := 0; i < maxPomoCount; i++ {
		if !end.After(e) {
			break
		}
		pomo := Pomo{0, "", Schedule{s, e}}
		c := false
		for _, mtg := range mtgs {
			if CheckConflictSchedule(pomo.Time, mtg.Time) {
				c = true
				s = mtg.Time.End
				e = s.Add(time.Duration(unit) * time.Minute)
				break
			}
		}
		if !c {
			count++
			pomo.Id = count
			pomos[count-1] = pomo
			s = e.Add(time.Duration(rest) * time.Minute)
			e = s.Add(time.Duration(unit) * time.Minute)
		}
	}
	return pomos
}

func display(mtgs []meeting, pomos []Pomo) {
	tasks := ReadTasks()
	cursor := 0
	point := "x"
	for _, pomo := range pomos {
		if len(mtgs) > 0 {
			for len(mtgs) > cursor && !mtgs[cursor].Time.Start.After(pomo.Time.Start) {
				displayMeeting(mtgs[cursor])
				cursor++
			}
		}
		if pomo.Id == 0 {
			break
		}
		if !time.Now().After(pomo.Time.End) && !time.Now().Before(pomo.Time.Start) {
			point = "*"
		} else if !time.Now().After(pomo.Time.End) {
			point = ""
		}
		task_str := buildTaskString(tasks, pomo.Id)
		displayPomo(pomo, point, task_str)
	}
	for i := cursor; i < len(mtgs); i++ {
		displayMeeting(mtgs[i])
	}
}

func buildTaskString(tasks []Task, id int) string {
	var task_arr []string
	for _, task := range tasks {
		if contains(task.Positions, id) {
			task_arr = append(task_arr, task.Title)
		}
	}
	return strings.Join(task_arr, ", ")
}

func displayPomo(pomo Pomo, point string, task_str string) {

	fmt.Printf("[%v] %v %s-%s %s %s\n", point, pomo.Id, pomo.Time.Start.Format(time_layout), pomo.Time.End.Format(time_layout), pomo.Name, task_str)
}

func displayMeeting(mtg meeting) {
	fmt.Printf("==== %s-%s %s =====\n", mtg.Time.Start.Format(time_layout), mtg.Time.End.Format(time_layout), mtg.Title)
}

func displayHeader() {
	const layout = "2006-01-02"
	fmt.Printf("#### %s\n\n", time.Now().Format(layout))
}
func displayFooter() {
	fmt.Println("-----")
	// day
	day := ReadDay()
	workingS := day.WorkingTime.Start
	workingE := day.WorkingTime.End
	fmt.Printf("working: %gh (%s-%s)\n", workingE.Sub(workingS).Hours(), workingS.Format(time_layout), workingE.Format(time_layout))
	// tasks
	tasks := ReadTasks()
	score := map[int]int{}
	fmt.Printf("task: ")
	for _, task := range tasks {
		for _, posi := range task.Positions {
			if v, ok := score[posi]; ok {
				score[posi] = v + 1
			} else {
				score[posi] = 1
			}
		}
	}
	for _, task := range tasks {
		var volume float64
		for _, posi := range task.Positions {
			v := float64(1) / float64(score[posi])
			volume = volume + v
		}
		mins := volume * float64(Conf.pomoTime+Conf.pomoRest)
		fmt.Printf("%s(%gh) ", task.Title, mins/60)
	}
	fmt.Printf("\n-----\n\n")
}

func init() {
	rootCmd.AddCommand(showCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	showCmd.Flags().BoolP("verbose", "v", false, "show header and footer")
}
