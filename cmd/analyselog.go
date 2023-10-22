package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/go-gota/gota/dataframe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ryannortham/digio-task/pkg/types"
)

// analyseLogCmd represents the parseLog command
var analyseLogCmd = &cobra.Command{
	Use:   "analyselog",
	Short: "Parses a log file containing HTTP requests and reports on its contents",

	Run: func(cmd *cobra.Command, args []string) {
		logReader := types.FileReader{Filename: getFilePath()}
		logParser := types.CombinedLogParser{}

		// read the log file
		logLines, err := logReader.ReadLines()
		if err != nil {
			fmt.Printf("Error reading log file: %v\n", err)
			return
		}

		// parse the log file
		var logEntries []types.LogEntry
		for _, line := range logLines {
			entry, err := logParser.ParseLogEntry(line)
			if err != nil {
				fmt.Printf("Error parsing log entry, omitting: %v\n", err)
				continue
			}

			logEntries = append(logEntries, entry)
		}

		// analyse the log file data
		df := dataframe.LoadStructs(logEntries)

		// count the number of unique IP addresses in the log
		uniqueIPs, err := getColumnSet(df, "IP")
		if err != nil {
			fmt.Printf("Error getting unique IPs: %v\n", err)
			return
		}

		fmt.Println("uniqueIPs:", len(uniqueIPs))
	},
}

func init() {
	rootCmd.AddCommand(analyseLogCmd)
	cobra.OnInitialize(initConfig)
}

func getFilePath() string {
	dir := viper.GetString("log-dir")
	fileName := viper.GetString("log-file")
	return filepath.Join(dir, fileName)
}

func getColumnSet(df dataframe.DataFrame, colName string) ([]interface{}, error) {
	//check if the column exists
	colNames := df.Names()

	for i, name := range colNames {
		if name == colName {
			break
		}
		if i == len(colNames)-1 {
			fmt.Printf("the column %s does not exist\n", colName)
			return nil, fmt.Errorf("invalid column name: %s", colName)
		}
	}

	// get the unique values in the column using a map
	var set []interface{}
	setMap := make(map[interface{}]bool)
	col := df.Col(colName)

	for _, val := range col.Records() {
		if _, ok := setMap[val]; !ok {
			setMap[val] = true
			set = append(set, val)
		}
	}

	return set, nil
}
