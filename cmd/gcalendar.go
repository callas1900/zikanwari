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
	Use:     "gcalendar",
	Aliases: []string{"gcal"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		main(cmd)
	},
}

func getCredentialFilePath() string {
	return GetConfigPath() + "credentials.json"
}

func getTokenFilePath() string {
	return GetConfigPath() + "token.json"
}

func main(cmd *cobra.Command) {
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
		const layout2 = "2006-01-02"
		displayDay := time.Now()
		for _, item := range events.Items {
			isDeclined := false
			for _, attendee := range item.Attendees {
				if attendee.Self {
					if attendee.ResponseStatus == "declined" {
						isDeclined = true
					}
					break
				}
			}
			if isDeclined {
				continue
			}

			noTimeSchedule := false
			start := item.Start.DateTime
			if start == "" {
				noTimeSchedule = true
				start = item.Start.Date
			}
			end := item.End.DateTime
			if end == "" {
				end = item.End.Date
			}
			var hours string
			var day time.Time
			if noTimeSchedule {
				day, _ = time.Parse(layout2, start)
			} else {
				hours, day = changeDateFormat(start, end)
			}
			if displayDay.YearDay() != day.YearDay() {
				fmt.Printf("\n%v:\n", day.Format(layout2))
				displayDay = day
			}

			fmt.Printf("[%v] %v\n", hours, item.Summary)
			if day.Day() == time.Now().Day() && day.Month() == time.Now().Month() && cmd.Flag("import").Changed && hours != "" {
				fmt.Println(hours)
				fmt.Println("---> imported")
				Set(item.Summary, hours)
			}
		}
	}
}

func changeDateFormat(start string, end string) (string, time.Time) {
	st, e := time.Parse(time.RFC3339Nano, start)
	if e != nil {
		fmt.Println()
	}
	et, e := time.Parse(time.RFC3339Nano, end)
	return st.Format(time_layout) + "-" + et.Format(time_layout), st
}

func init() {
	rootCmd.AddCommand(gcalendarCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gcalendarCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	gcalendarCmd.Flags().BoolP("import", "i", false, "import from gcalendar")
}
