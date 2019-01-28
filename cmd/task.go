// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "A brief description of your command",
	Long: `create tasks.
For example:
	* zikanwari task do something 1,2,3
	* zikanwari task read email 3`,
	Run: func(cmd *cobra.Command, args []string) {
		// posi_str := strings.Split(args[len(args)-1], ",")
		posi, _ := strconv.Atoi(args[len(args)-1])
		args = args[:len(args)-1]
		title := strings.Join(args, " ")
		AddTask(title, posi)
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// taskCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// taskCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
