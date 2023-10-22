package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ryannortham/digio-task/log"
)

// analyseLogCmd represents the parseLog command
var analyseLogCmd = &cobra.Command{
	Use:   "analyselog",
	Short: "Parses a log file containing HTTP requests and reports on its contents",

	Run: func(cmd *cobra.Command, args []string) {
		logReader := &log.FileReader{Filename: getFilePath()}
		logParser := &log.CombinedLogParser{}
		logAnalyzer := &log.LogAnalyzer{}

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
		logAnalysis, err := logAnalyzer.AnalyseLogEntries(logEntries)
		if err != nil {
			fmt.Printf("Error analysing log file: %v\n", err)
			return
		}

		fmt.Printf("unique ip addresses: %d\n", logAnalysis.UniqueIPCount)

	},
}

func init() {
	rootCmd.AddCommand(analyseLogCmd)
}

func getFilePath() string {
	dir := viper.GetString("log-dir")
	fileName := viper.GetString("log-file")
	return filepath.Join(dir, fileName)
}
