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

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "read ,write and show list config file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			fmt.Println("config file : " + viper.ConfigFileUsed())
			fmt.Println("---")
			for _, key := range viper.AllKeys() {
				fmt.Printf("%s : %v\n", key, viper.Get(key))
			}
			break
		case 1:
			key := args[0]
			if viper.IsSet(key) {
				fmt.Printf("%v\n", viper.Get(key))
			} else {
				fmt.Println("key not found")
			}
			break
		case 2:
			key, val := args[0], args[1]
			if viper.IsSet(key) {
				preval := viper.Get(key)
				viper.Set(key, val)
				viper.WriteConfig()
				fmt.Printf("update! %s %v -> %v\n", key, preval, val)
			} else {
				fmt.Println("key not found")
			}
			break
		default:
			fmt.Println("error")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
