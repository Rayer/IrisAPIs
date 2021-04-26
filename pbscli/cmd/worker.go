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
	"context"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"time"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		for {
			data, _ := pbsService.FetchPbsFromServer(context.TODO())
			//bar := progressbar.Default(int64(len(data)))
			bar := progressbar.NewOptions(len(data),
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionSetWidth(15),
				progressbar.OptionSetDescription("[cyan][1/3][reset] Writing moshable file..."),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "[green]=[reset]",
					SaucerHead:    "[green]>[reset]",
					SaucerPadding: " ",
					BarStart:      "[",
					BarEnd:        "]",
				}))

			err := pbsService.UpdateDatabase(context.TODO(), data, func(total int, now int, updated int, inserted int, skipped int) {
				bar.Describe(fmt.Sprintf("[cyan](%4d/%4d)[reset] %4d updated, %4d inserted %4d skipped", now, total, updated, inserted, skipped))
				bar.Add(1)
			})
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println()
			time.Sleep(1 * time.Minute)
		}
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
