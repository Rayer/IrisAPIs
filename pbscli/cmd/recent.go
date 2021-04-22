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
	"os"
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
		hours, err := cmd.Flags().GetInt("time")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res, err := s.GetHistory(context.TODO(), time.Duration(hours)*time.Hour)
		if err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range res {
			fmt.Println("ID : ", k)
			for _, events := range v {
				fmt.Printf("%s\t%s\n", events.LastUpdateTimestamp.Format(time.Stamp), *events.Information)
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(recentCmd)
	recentCmd.Flags().IntP("time", "t", 6, "Recent n hours")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
