/*
Copyright © 2020 Rayer

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
	"fmt"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "apikey_cli",
	Short: "An utility to manage api keys for IrisAPIs",
	Long:  "An utility to manage api keys : list, update, issue and expire",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var cfgFile string
var connectionString string
var service IrisAPIs.ApiKeyService
var verbose bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config/iris-apis.yaml and ./iris-apis.yaml)")
	rootCmd.PersistentFlags().StringVar(&connectionString, "connection-string", "", "Connection string to database")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	_ = viper.BindPFlag("ConnectionString", rootCmd.PersistentFlags().Lookup("connection-string"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
		viper.SetConfigName("iris-apis")
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	rootCmd.PersistentFlags().HasFlags()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		//dbContext, err = IrisAPIs.NewDatabaseContext(viper.GetString("ConnectionString"), false)
		//service = IrisAPIs.NewApiKeyService(dbContext)
	}
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	dbContext, err := IrisAPIs.NewDatabaseContext(viper.GetString("ConnectionString"), verbose)
	if err != nil {
		panic(err)
	}
	service = IrisAPIs.NewApiKeyService(dbContext)

}
