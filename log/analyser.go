package log

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
)

type LogAnalysis struct {
	UniqueIPCount       int
	TopNMostVisitedURLs [][]string
	TopNMostActiveIPs   [][]string
}

type LogAnalyzer interface {
	GetLogAnalysis([]LogEntry, int) (*LogAnalysis, error)
}

type CombinedLogAnalyzer struct{}

func (l *CombinedLogAnalyzer) GetLogAnalysis(logEntries []LogEntry, topN int) (*LogAnalysis, error) {
	df := dataframe.LoadStructs(logEntries)
	if df.Err != nil {
		return nil, df.Err
	}

	IPGroups, err := aggregateDfByColumn(df, "IP")
	if err != nil {
		return nil, err
	}

	URLGroups, err := aggregateDfByColumn(df, "URL")
	if err != nil {
		return nil, err
	}

	topActiveIPs, err := getTopNRows(IPGroups, topN)
	if err != nil {
		return nil, err
	}

	topVisitedURLs, err := getTopNRows(URLGroups, topN)
	if err != nil {
		return nil, err
	}

	la := &LogAnalysis{
		UniqueIPCount:       IPGroups.Nrow(),
		TopNMostActiveIPs:   topActiveIPs.Records(),
		TopNMostVisitedURLs: topVisitedURLs.Records(),
	}

	return la, nil
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

	// sort the dataframe to make output deterministic
	groupedDf = groupedDf.Arrange(dataframe.Sort(colName))

	return &groupedDf, nil
}

func getTopNRows(df *dataframe.DataFrame, n int) (dataframe.DataFrame, error) {
	if n > df.Nrow() {
		return dataframe.DataFrame{}, fmt.Errorf("n is greater than the number of rows in the dataframe")
	}

	// sort the dataframe by the count column then name column
	sortedDf := df.Arrange(dataframe.Sort(df.Names()[0])).Arrange(dataframe.RevSort(df.Names()[1]))

	indices := make([]int, n)
	for i := range indices {
		// skip the first row as it is the header
		if i == 0 {
			continue
		}

		indices[i] = i
	}

	return sortedDf.Subset(indices), nil
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
