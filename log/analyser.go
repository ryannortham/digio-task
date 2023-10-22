package log

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
)

type LogAnalyzer struct{}

type LogAnalysis struct {
	UniqueIPCount int
	// top3MostVisitedURLs []string
	// top3MostActiveIPs   []string
}

func (l *LogAnalyzer) AnalyseLogEntries(logEntries []LogEntry) (LogAnalysis, error) {
	df := dataframe.LoadStructs(logEntries)

	// count the number of unique IP addresses in the log
	uniqueIPs, err := l.getColumnSet(df, "IP")
	if err != nil {
		return LogAnalysis{}, fmt.Errorf("error getting unique IPs: %w", err)
	}

	result := LogAnalysis{
		UniqueIPCount: len(uniqueIPs),
	}

	return result, err
}

// getColumnSet returns a slice of unique values in a column
func (l *LogAnalyzer) getColumnSet(df dataframe.DataFrame, colName string) ([]interface{}, error) {

	// check that the column exists
	if !columnExists(df, colName) {
		return nil, fmt.Errorf("column %s does not exist", colName)
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

// check that a column exists in a dataframe
func columnExists(df dataframe.DataFrame, colName string) bool {
	colNames := df.Names()

	for _, name := range colNames {
		if name == colName {
			return true
		}
	}

	return false
}
