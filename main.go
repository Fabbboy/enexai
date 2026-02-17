package main

import (
	"github.com/spf13/cobra"
	"swisscom.com/fabrice/enexai/impl"
)

var configPath string
var csvPath string

const (
	configLabel   string = "config"
	configDefault string = "config.ini"
	csvLabel      string = "csv"
	csvDefault    string = "skills.xlsx"
)

var rootCmd = &cobra.Command{
	Use: "app",
	Run: func(cmd *cobra.Command, args []string) {
		impl.Run(configPath, csvPath)
	},
}

func init() {
	rootCmd.Flags().StringVar(&configPath, configLabel, configDefault, "")
	rootCmd.Flags().StringVar(&csvPath, csvLabel, csvDefault, "")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
