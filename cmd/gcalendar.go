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
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// gcalendarCmd represents the gcalendar command
var gcalendarCmd = &cobra.Command{
	Use:   "gcalendar",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

func getCredentialFilePath() string {
	return GetConfigPath() + "credentials.json"
}

func getTokenFilePath() string {
	return GetConfigPath() + "token.json"
}

func main() {
	b, err := ioutil.ReadFile(getCredentialFilePath())
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		const layout = "01-02 15:04"
		for _, item := range events.Items {
			start := item.Start.DateTime
			if start == "" {
				start = item.Start.Date
			}
			end := item.End.DateTime
			if end == "" {
				end = item.End.Date
			}
			fmt.Printf("[%v-%v] %v)\n", changeDateFormat(start, layout), changeDateFormat(end, time_layout), item.Summary)
		}
	}
}

func changeDateFormat(in string, layout string) string {
	t, _ := time.Parse(time.RFC3339Nano, in)
	return t.Format(layout)
}

func init() {
	rootCmd.AddCommand(gcalendarCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gcalendarCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gcalendarCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
