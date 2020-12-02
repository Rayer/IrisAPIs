/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

// usageCmd represents the usage command
var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Show usage for specified API Key",
	Long:  "Show usage statistics for specified API Key, will aggregated into IPs, Paths and Occurrence",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 || func() bool {
			_, err := strconv.Atoi(args[0])
			return err != nil
		}() {
			return errors.New("require exactly 1 id as argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("usage called, args is : ", args)
		v, _ := strconv.Atoi(args[0])
		days, _ := cmd.Flags().GetInt("days")
		now := time.Now()
		before := time.Now().AddDate(0, 0, -days)
		res, err := service.GetKeyUsage(v, &before, &now)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, u := range res {
			fmt.Printf("key id : %d, path : %s, timestamp : %s\n", *u.ApiKeyRef, *u.Fullpath, u.Timestamp.Format(time.Stamp))
		}
	},
}

func init() {
	rootCmd.AddCommand(usageCmd)
	usageCmd.Flags().IntP("days", "d", 7, "Define usage in n days.")
	//usageCmd.Flags().BoolP("byPath", "p", false, "Use path as argument instead of Key ID")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
