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
	"time"

	"github.com/spf13/cobra"

	"encoding/json"
	"io/ioutil"
	"log"
)

type Pomo struct {
	Id        int       `json:"id"`
	Name      string    `json:"title"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
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
		content, err := ioutil.ReadFile(Conf.DataJSONPath)
		if err != nil {
			log.Fatal(err)
		}
		var mtg meeting
		json.Unmarshal(content, &mtg)
		// TODO: dummy impl
		now := time.Now()
		startD := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
		endD := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
		pomos := CalcPomos([]meeting{mtg}, startD, endD, 25, 5)
		display([]meeting{mtg}, pomos)
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
		pomo := Pomo{0, "", s, e}
		for _, mtg := range mtgs {
			if checkConflict(mtg, pomo) {
				count++
				pomo.Id = count
				pomos[count-1] = pomo
			}
		}
		s = e.Add(time.Duration(rest) * time.Minute)
		e = s.Add(time.Duration(unit) * time.Minute)
	}
	return pomos
}

func checkConflict(mtg meeting, pomo Pomo) bool {
	return !pomo.EndTime.After(mtg.StartTime) || !pomo.StartTime.Before(mtg.EndTime)
}

func display(mtgs []meeting, pomos []Pomo) {
	const layout = "15:04"
	for _, pomo := range pomos {
		if pomo.Id == 0 {
			break
		}
		for _, mtg := range mtgs {
			if !mtg.StartTime.After(pomo.StartTime) && !mtg.EndTime.Before(pomo.StartTime) {
				fmt.Printf("==== %s %s-%s =====\n", mtg.Title, mtg.StartTime.Format(layout), mtg.EndTime.Format(layout))
			}
		}
		fmt.Printf("[] %v %s-%s %s \n", pomo.Id, pomo.StartTime.Format(layout), pomo.EndTime.Format(layout), pomo.Name)
	}
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
