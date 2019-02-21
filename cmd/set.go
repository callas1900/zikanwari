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
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set meeting schedule",
	Long: `set meeting schedule.
For example:
	* zikanwari set Mtg 11:00-13:00
	* zikanwari set Launch with team 11:00-13:00
`,
	Run: func(cmd *cobra.Command, args []string) {
		schedule := BuildScheduleStruct(args[len(args)-1])
		if !checkMeetingAvailable(schedule.Start, schedule.End) {
			fmt.Println("conflict with another meeting")
			os.Exit(1)
		}
		args = args[:len(args)-1]
		title := strings.Join(args, " ")
		m := meeting{1, title, schedule}
		AddMeeting(m)
	},
}

func checkMeetingAvailable(start time.Time, end time.Time) bool {
	mtgs := ReadMeetings()
	source := Schedule{start, end}
	for _, mtg := range mtgs {
		check := CheckConflictSchedule(source, mtg.Time)
		if check {
			return true
		}
	}
	return true
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
