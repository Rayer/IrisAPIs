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
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strconv"
)

// expireCmd represents the expire command
var expireCmd = &cobra.Command{
	Use:   "expire",
	Short: "Expire an API Key or verse visa.",
	Long:  "Set an API Key to expire, or cancel expiration",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("require argument as api key ID")
		}
		_, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("expire called")
		id, _ := strconv.Atoi(args[0])
		enable, _ := cmd.Flags().GetBool("re-enable")
		err := service.SetExpire(id, !enable)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(expireCmd)
	expireCmd.Flags().BoolP("re-enable", "r", false, "Re-enable expired api key.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// expireCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// expireCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
