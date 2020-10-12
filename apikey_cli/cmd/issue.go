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
	"github.com/spf13/cobra"
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issue [application_name]",
	Short: "Issue an new API Key",
	Long:  `Issue an new API Key. API Key is used in some limited API endpoints, and moreover, some sensitive endpoints can be protected by only accessible via Privileged API Keys`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f := cmd.Flags()
		application := args[0]
		privileged, _ := f.GetBool("privileged")
		issuer, _ := f.GetString("issuer")

		key, err := service.IssueApiKey(application, true, true, issuer, privileged)
		if err != nil {
			panic(err)
		}
		fmt.Println("API Key : ", key)
	},
}

func init() {
	rootCmd.AddCommand(issueCmd)

	flags := issueCmd.Flags()
	//flags.StringP("application", "a", "", "Which application will use this API Key (required)")
	flags.StringP("issuer", "i", "auto", "Issuer")
	flags.BoolP("privileged", "p", false, "Privileged API Key")

	//issueCmd.MarkFlagRequired("application")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// issueCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// issueCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
