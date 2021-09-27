/*
Package cmd
Copyright © 2021 Rayer Tung rayer@vista.aero

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
	"github.com/spf13/cobra"
)

// singlestepCmd represents the singlestep command
var singlestepCmd = &cobra.Command{
	Use:   "singlestep",
	Short: "Start Iris data fetcher",
	Long:  `Start Iris data fetcher. It is suppliant application for Iris API, gathering data into page`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return StartTaskLoop()
	},
}

func init() {
	rootCmd.AddCommand(singlestepCmd)
}
