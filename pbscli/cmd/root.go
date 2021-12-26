// Package cmd /*
package cmd

import (
	"context"
	"fmt"
	"github.com/Rayer/IrisAPIs"
	"github.com/spf13/cobra"
	"github.com/xormplus/xorm/log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var dbContext *IrisAPIs.DatabaseContext
var pbsService IrisAPIs.PbsTrafficDataService
var ctx context.Context

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pbs",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	ctx = context.Background()
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pbs.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	dbContext, err := IrisAPIs.NewTestDatabaseContext(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	dbContext.DbObject.Logger().ShowSQL(false, false, false)
	dbContext.DbObject.Logger().SetLevel(log.LOG_WARNING)
	pbsService = IrisAPIs.NewPbsTrafficDataService(dbContext)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pbs" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pbs")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
