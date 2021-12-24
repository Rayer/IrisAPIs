// Package cmd /*
package cmd

import (
	"context"
	"fmt"
	"github.com/Rayer/IrisAPIs"
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
		s := IrisAPIs.NewPbsTrafficDataService(dbContext)
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
