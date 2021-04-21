/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"IrisAPIs"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

// recentCmd represents the recent command
var recentCmd = &cobra.Command{
	Use:   "recent",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		s := IrisAPIs.NewPbsTrafficDataService(dbContext.DbObject)
		res, err := s.GetHistory(context.TODO(), 12*time.Hour)
		if err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range res {
			fmt.Println("ID : ", k)
			for i, events := range v {
				if i == 0 {
					fmt.Printf("%s\t%s\n", events.EntryTimestamp.Format(time.Stamp), *events.CurInfo)
				}
				if events.HistoryInfo != nil {
					fmt.Printf("%s\t%s\n", events.Timestamp.Format(time.Stamp), *events.HistoryInfo)
				}
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(recentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
