/*
Package cmd
Copyright © 2021 Rayer Tung

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"IrisAPIs"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//var cfgFile string
//var config *IrisAPIs.Configuration
var connStr string
var fixioKey string
var successPeriod int
var failPeriod int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "datafetch",
	Short: "Data fetch service for Iris Mainframe APIs",
	Long:  `Data fetch service, it will update data that Iris APIs will used.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.iris-datafetch.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	singlestepCmd.PersistentFlags().StringVar(&connStr, "connection_string", "", "Connection String to SQL")
	singlestepCmd.PersistentFlags().StringVar(&fixioKey, "fixio_key", "", "Fixio API Key")
	singlestepCmd.PersistentFlags().IntVar(&successPeriod, "success_period", 43200, "Fetch frequency after Success")
	singlestepCmd.PersistentFlags().IntVar(&failPeriod, "fail_period", 10800, "Next fetch attempt after fail")
}

//// initConfig reads in config file and ENV variables if set.
//func initConfig() {
//	if cfgFile != "" {
//		// Use config file from the flag.
//		viper.SetConfigFile(cfgFile)
//	} else {
//		// Find home directory.
//		home, err := homedir.Dir()
//		cobra.CheckErr(err)
//
//		// Search config in home directory with name ".iris-datafetch" (without extension).
//		viper.AddConfigPath(home)
//		viper.SetConfigName("iris-apis")
//		viper.SetConfigType("yaml")
//	}
//
//	viper.AutomaticEnv() // read in environment variables that match
//	_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
//	config = IrisAPIs.NewConfiguration()
//}

func StartTaskLoop() error {

	//Validate arguments
	if connStr == "" {
		return errors.New("\"--connection_string\" is required!")
	}

	if fixioKey == "" {
		return errors.New("\"--fixio_key\" is required!")
	}

	fmt.Println("singlestep called")
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT)
	ctx := context.Background()
	database, _ := IrisAPIs.NewDatabaseContext(connStr, true, nil)
	service := IrisAPIs.NewCurrencyContextWithConfig(fixioKey, successPeriod, failPeriod, database)
	service.CurrencySyncRoutine(ctx)
	timerSuccess := time.NewTimer(time.Duration(successPeriod) * time.Second)
	timerFail := time.NewTicker(time.Duration(failPeriod) * time.Second)
	timerFail.Stop()
OUTERLOOP:
	for {
		select {
		case sig := <-sigc:
			fmt.Println("Intercepted signal", sig)
			break OUTERLOOP
		case <-timerSuccess.C:
		case <-timerFail.C:
			result, err := service.CurrencySyncWorker()
			if err != nil {
				timerFail.Reset(time.Duration(failPeriod) * time.Second)
			} else {
				timerFail.Stop()
			}
			fmt.Println(result)
		default:
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}
