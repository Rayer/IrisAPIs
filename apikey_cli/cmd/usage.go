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
	"context"
	"errors"
	"fmt"
	"github.com/Rayer/IrisAPIs"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// usageCmd represents the usage command
var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Show usage for specified API Key",
	Long:  "Show usage statistics for specified API Key, will aggregated into IPs, Paths and Occurrence",
	Args: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return errors.New("require exactly 1 argument")
		}

		if func() bool {
			b, _ := cmd.Flags().GetBool("byPath")
			return b
		}() {
			return nil
		} else if func() bool {
			_, err := strconv.Atoi(args[0])
			return err != nil
		}() {
			return errors.New("error parsing argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		days, _ := cmd.Flags().GetInt("days")
		now := time.Now()
		prev := time.Now().AddDate(0, 0, -days)
		var retValue []*IrisAPIs.ApiKeyAccess
		if func() bool {
			b, _ := cmd.Flags().GetBool("byPath")
			return b
		}() {
			retValue, _ = service.GetKeyUsageByPath(context.TODO(), args[0], false, &prev, &now)
		} else {
			retValue, _ = service.GetKeyUsageById(context.TODO(), func() int {
				i, _ := strconv.Atoi(args[0])
				return i
			}(), &prev, &now)
		}

		for _, v := range retValue {
			fmt.Printf("Key id : %d, path : %s(%s), ip : %s(%s), time : %s\n", *v.ApiKeyRef, *v.Fullpath, *v.Method, *v.Ip, *v.Nation, v.Timestamp.Format(time.RFC822))
		}
	},
}

func init() {
	rootCmd.AddCommand(usageCmd)
	usageCmd.Flags().BoolP("byPath", "p", false, "Use this flag to mark filter by path")
	usageCmd.Flags().IntP("days", "d", 7, "Retrieve n days")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
