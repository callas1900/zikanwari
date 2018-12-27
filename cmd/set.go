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
	"strings"
	"time"
	"encoding/json"
	"io/ioutil"

	"github.com/spf13/cobra"
)

type meeting struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set meeting schedule",
	Long: `set meeting schedule.
For example:

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		inputs := strings.Split(args[1], "-")
		start, end := buildInputDays(inputs)
		m := meeting{1, args[0], start, end}
		b, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(m)
		ioutil.WriteFile("data.json", b, 0644)
	},
}

func buildInputDays(days []string) (string, string) {
	const layout = "2006-01-02 15:04"
	prefix := strings.Split(time.Now().Format(layout), " ")[0]
	s := strings.Join([]string{prefix, days[0]}, " ")
	e := strings.Join([]string{prefix, days[1]}, " ")
	return s, e
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
