package log

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
)

type LogAnalysis struct {
	UniqueIPCount       int
	Top3MostVisitedURLs [][]string
	Top3MostActiveIPs   [][]string
}

type LogAnalyzer interface {
	GetLogAnalysis([]LogEntry) (*LogAnalysis, error)
}

type CombinedLogAnalyzer struct{}

// GetLogAnalysis returns a LogAnalysis struct containing the results of the analysis
func (l *CombinedLogAnalyzer) GetLogAnalysis(logEntries []LogEntry) (*LogAnalysis, error) {
	df := dataframe.LoadStructs(logEntries)

	IPGroups, err := aggregateDfByColumn(df, "IP")
	if err != nil {
		return nil, err
	}

	URLGroups, err := aggregateDfByColumn(df, "URL")
	if err != nil {
		return nil, err
	}

	la := &LogAnalysis{
		UniqueIPCount:       IPGroups.Nrow(),
		Top3MostActiveIPs:   getTopNRows(IPGroups, 3).Records(),
		Top3MostVisitedURLs: getTopNRows(URLGroups, 3).Records(),
	}

	return la, err
}

func aggregateDfByColumn(df dataframe.DataFrame, colName string) (*dataframe.DataFrame, error) {
	if !columnExists(df, colName) {
		return nil, fmt.Errorf("column %s does not exist", colName)
	}

	// group the dataframe by the column and aggregate the count of each group
	groupedDf := df.GroupBy(colName).Aggregation(
		[]dataframe.AggregationType{dataframe.Aggregation_COUNT},
		[]string{colName},
	)

	return &groupedDf, nil
}

func getTopNRows(df *dataframe.DataFrame, n int) dataframe.DataFrame {
	// sort the dataframe by the count column
	sortedDf := df.Arrange(dataframe.RevSort(df.Names()[1]))

	indices := make([]int, n)
	for i := range indices {
		// skip the first row as it is the header
		if i == 0 {
			continue
		}

		indices[i] = i
	}

	return sortedDf.Subset(indices)
}

// checks that a column exists in a dataframe
func columnExists(df dataframe.DataFrame, colName string) bool {
	for _, name := range df.Names() {
		if name == colName {
			return true
		}
	}
	return false
}
