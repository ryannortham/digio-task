package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ryannortham/digio-task/log"
	"github.com/ryannortham/digio-task/ui"
)

var (
	logReader   log.LogReader
	logParser   log.LogParser
	logAnalyzer log.LogAnalyzer

	rootCmd = &cobra.Command{
		Short: "Parses a log file containing HTTP requests and to reports on its contents",
		Long: `
Parses a log file containing HTTP requests and to reports on its contents

For a given log file we want to know:
- The number of unique IP addresses
- The top 3 most visited URLs
- The top 3 most active IP addresses
`,

		Run: func(cmd *cobra.Command, args []string) {
			// read the log file
			logLines, err := logReader.ReadLines()
			if err != nil {
				fmt.Printf("Error reading log file: %v\n", err)
				return
			}

			// parse the log file
			logEntries, err := logParser.ParseLogEntries(logLines)
			if err != nil {
				fmt.Printf("Error parsing log file: %v\n", err)
				return
			}

			// analyse the log file data
			logAnalysis, err := logAnalyzer.GetLogAnalysis(logEntries)
			if err != nil {
				fmt.Printf("Error analysing log file: %v\n", err)
				return
			}

			// print the results
			ui.PrintAnalysisResults(logAnalysis)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLogReader)
	cobra.OnInitialize(initLogParser)
	cobra.OnInitialize(initLogAnalyzer)
}

func initConfig() {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s", err)
	}
}

func initLogReader() {
	logSource := viper.GetString("log-source")

	switch logSource {
	case "file":
		logFilePath := filepath.Join(viper.GetString("log-dir"), viper.GetString("log-file"))
		logReader = &log.FileReader{LogFilePath: logFilePath}
	case "api":
		fmt.Println("API log source not yet implemented")
		os.Exit(1)
	default:
		fmt.Printf("Unknown log source: %s\n", logSource)
		os.Exit(1)
	}
}

func initLogParser() {
	logFormat := viper.GetString("log-format")

	switch logFormat {
	case "combined-log-format":
		logParser = &log.CombinedLogParser{}
	case "common-log-format":
		fmt.Println("Common log format not yet implemented")
		os.Exit(1)
	default:
		fmt.Printf("Unknown log format: %s\n", logFormat)
		os.Exit(1)
	}
}

func initLogAnalyzer() {
	logFormat := viper.GetString("log-format")

	switch logFormat {
	case "combined-log-format":
		logAnalyzer = &log.CombinedLogAnalyzer{}
	case "common-log-format":
		fmt.Println("Common log format not yet implemented")
		os.Exit(1)
	default:
		fmt.Printf("Unknown log format: %s\n", logFormat)
		os.Exit(1)
	}
}
