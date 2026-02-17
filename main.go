package main

import (
	"github.com/fabrice/enexai/impl"
	"github.com/spf13/cobra"
)

var configPath string
var csvPath string
var debug bool

const (
	configLabel   string = "config"
	configDefault string = "config.ini"
	csvLabel      string = "csv"
	csvDefault    string = "skills.csv"
	debugLabel    string = "debug"
)

var rootCmd = &cobra.Command{
	Use: "app",
	Run: func(cmd *cobra.Command, args []string) {
		err := impl.Run(configPath, csvPath, debug)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&configPath, configLabel, configDefault, "")
	rootCmd.Flags().StringVar(&csvPath, csvLabel, csvDefault, "")
	rootCmd.Flags().BoolVar(&debug, debugLabel, false, "enable debug logging")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
