// Package cmd /*
package cmd

import (
	"context"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"time"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update PBS entries",
	Long: `Update PBS from server, and write into database. 
It can be run either in one-time mode(without argument), worker mode(with -w argument) and daemon mode(-d daemon)`,
	Run: func(cmd *cobra.Command, args []string) {
		isWorker, _ := cmd.Flags().GetBool("worker")
		isDaemon, _ := cmd.Flags().GetBool("daemon")
		updateTime, _ := cmd.Flags().GetInt("update_timer")
		if isDaemon {
			//NYI
		} else {
			fetchAndPrint()
			if isWorker {
				for {
					t := time.NewTicker(time.Duration(updateTime) * time.Minute)
					select {
					case <-t.C:
						fetchAndPrint()
					}
				}
			}
		}
	},
}

func fetchAndPrint() {
	data, _ := pbsService.FetchPbsFromServer(context.TODO())
	bar := progressbar.NewOptions(len(data),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
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
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().IntP("update_timer", "u", 1, "Continuous update every n min. It only work with -w(worker) or -d (daemon) command")
	updateCmd.Flags().BoolP("worker", "w", false, "Worker mode")
	updateCmd.Flags().BoolP("daemon", "d", false, "Daemon mode(TBD)")
}
