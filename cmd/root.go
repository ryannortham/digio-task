package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ryannortham/digio-task/log"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "Parses a log file containing HTTP requests and to reports on its contents",
	Long: `
Parses a log file containing HTTP requests and to reports on its contents

For a given log file we want to know:
  - The number of unique IP addresses
  - The top 3 most visited URLs
  - The top 3 most active IP addresses
`,

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

		fmt.Print(getDigioLogo())
		fmt.Printf("unique ip addresses: %d\n", logAnalysis.UniqueIPCount)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s", err)
	}

	fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
}

func getFilePath() string {
	dir := viper.GetString("log-dir")
	fileName := viper.GetString("log-file")
	return filepath.Join(dir, fileName)
}

func getDigioLogo() string {
	return `
    
         xxxxxx                                    $$$$$$   $$$$$$                          $$$$$$                      
         xxxxxx       :                            $$$$$$   $$$$$$                          $$$$$$                      
        xxxxxx    :::::                            $$$$$$   $$$$$$                          $$$$$$                      
      xxxxxxxx  ::::::::                           $$$$$$                                                               
  xxxxxxxxxxx ::::::::::   ++            $$$$$$$$$ $$$$$$   $$$$$$      $$$$$$$$$$ $$$$$$   $$$$$$       $$$$$$$$$$     
xxxxxxxxxxx  ::::::::    +++++         $$$$$$$$$$$$$$$$$$   $$$$$$    $$$$$$$$$$$$$$$$$$$   $$$$$$    $$$$$$$$$$$$$$$$  
xxxxxxxxx    ::::::     +++++++      $$$$$$$$$$$$$$$$$$$$   $$$$$$   $$$$$$$$$$$$$$$$$$$$   $$$$$$   $$$$$$$$$$$$$$$$$$ 
xxxxxx      ::::::     ++++++++      $$$$$$       $$$$$$$   $$$$$$  $$$$$$$      $$$$$$$$   $$$$$$  $$$$$$$      $$$$$$$
            ::::::    +++++++       $$$$$$         $$$$$$   $$$$$$  $$$$$$         $$$$$$   $$$$$$  $$$$$$        $$$$$$
   ;;;      ::::::   +++++++        $$$$$$         $$$$$$   $$$$$$  $$$$$$         $$$$$$   $$$$$$  $$$$$          $$$$$
 ;;;;;;     ::::::   ++++++         $$$$$$         $$$$$$   $$$$$$  $$$$$$        $$$$$$$   $$$$$$  $$$$$$        $$$$$$
 ;;;;;;;;   ::::     ++++++          $$$$$$       $$$$$$$   $$$$$$   $$$$$$$$$ $$$$$$$$$$   $$$$$$  $$$$$$$      $$$$$$$
  ;;;;;;;;;          ++++++          $$$$$$$$$$$$$$$$$$$$   $$$$$$   $$$$$$$$$$$$$$$$$$$$   $$$$$$   $$$$$$$$$$$$$$$$$$ 
    ;;;;;;;;;;;;;;   +++++++           $$$$$$$$$$$$$$$$$$   $$$$$$     $$$$$$$$$$$$$$$$$$   $$$$$$    $$$$$$$$$$$$$$$$  
     ;;;;;;;;;;;;;;   +++++              $$$$$$$$$ $$$$$$   $$$$$$        $$$$$$   $$$$$$   $$$$$$       $$$$$$$$$$$    
        ;;;;;;;;;;;    +                                                $$$       $$$$$$                                
                                                                      $$$$$$$$$$$$$$$$$$                                
                                                                       $$$$$$$$$$$$$$$$                                 
                                                                         $$$$$$$$$$$$         

`
}
