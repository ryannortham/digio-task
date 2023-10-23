package log

import (
	"testing"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/stretchr/testify/assert"
)

func Test_columnExists(t *testing.T) {
	df := dataframe.New(
		series.New([]int{1, 2, 3}, series.Int, "A"),
		series.New([]int{4, 5, 6}, series.Int, "B"),
	)

	tests := []struct {
		name    string
		df      dataframe.DataFrame
		colName string
		want    bool
	}{
		{
			name:    "column exists",
			df:      df,
			colName: "B",
			want:    true,
		},
		{
			name:    "column does not exist",
			df:      df,
			colName: "C",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := columnExists(tt.df, tt.colName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getTopNRows(t *testing.T) {
	df := dataframe.New(
		series.New([]string{"Alice", "Bob", "Charlie", "David"}, series.String, "Name"),
		series.New([]int{10, 20, 30, 40}, series.Int, "Count"),
	)

	tests := []struct {
		name    string
		df      *dataframe.DataFrame
		n       int
		want    dataframe.DataFrame
		wantErr bool
	}{
		{
			name: "get top 2 rows",
			df:   &df,
			n:    2,
			want: dataframe.New(
				series.New([]string{"David", "Charlie"}, series.String, "Name"),
				series.New([]int{40, 30}, series.Int, "Count"),
			),
			wantErr: false,
		},
		{
			name: "get all rows",
			df:   &df,
			n:    4,
			want: dataframe.New(
				series.New([]string{"David", "Charlie", "Bob", "Alice"}, series.String, "Name"),
				series.New([]int{40, 30, 20, 10}, series.Int, "Count"),
			),
			wantErr: false,
		},
		{
			name: "get 0 rows",
			df:   &df,
			n:    0,
			want: dataframe.New(
				series.New([]string{}, series.String, "Name"),
				series.New([]int{}, series.Int, "Count"),
			),
			wantErr: false,
		},
		{
			name:    "get more rows than dataframe has",
			df:      &df,
			n:       5,
			want:    dataframe.DataFrame{},
			wantErr: true,
		},
		{
			name:    "get top N rows of empty dataframe",
			df:      &dataframe.DataFrame{},
			n:       2,
			want:    dataframe.DataFrame{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTopNRows(tt.df, tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTopNRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_aggregateDfByColumn(t *testing.T) {
	df := dataframe.New(
		series.New([]string{"Alice", "Bob", "Charlie", "David", "Alice"}, series.String, "Name"),
		series.New([]int{10, 20, 30, 40, 50}, series.Int, "Count"),
	)

	tests := []struct {
		name    string
		df      dataframe.DataFrame
		colName string
		want    dataframe.DataFrame
		wantErr bool
	}{
		{
			name:    "aggregate by existing column",
			df:      df,
			colName: "Name",
			want: dataframe.New(
				series.New([]string{"Alice", "Bob", "Charlie", "David"}, series.String, "Name"),
				series.New([]float64{2, 1, 1, 1}, series.Float, "Name_COUNT"),
			),
			wantErr: false,
		},
		{
			name:    "aggregate by non-existing column",
			df:      df,
			colName: "Age",
			want:    dataframe.DataFrame{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := aggregateDfByColumn(tt.df, tt.colName)
			if (err != nil) != tt.wantErr {
				t.Errorf("aggregateDfByColumn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				assert.Equal(t, tt.want, *got)
			}
		})
	}
}

func Test_CombinedLogAnalyzer_GetLogAnalysis(t *testing.T) {
	logEntries := []LogEntry{
		{IP: "192.168.0.1", URL: "/home"},
		{IP: "192.168.0.2", URL: "/about"},
		{IP: "192.168.0.1", URL: "/home"},
		{IP: "192.168.0.3", URL: "/contact"},
		{IP: "192.168.0.1", URL: "/home"},
		{IP: "192.168.0.2", URL: "/about"},
	}

	tests := []struct {
		name    string
		entries []LogEntry
		topN    int
		want    *LogAnalysis
		wantErr bool
	}{
		{
			name:    "get log analysis with valid input is successful",
			entries: logEntries,
			topN:    2,
			want: &LogAnalysis{
				UniqueIPCount:       3,
				TopNMostActiveIPs:   [][]string{{"IP", "IP_COUNT"}, {"192.168.0.1", "3.000000"}, {"192.168.0.2", "2.000000"}},
				TopNMostVisitedURLs: [][]string{{"URL", "URL_COUNT"}, {"/home", "3.000000"}, {"/about", "2.000000"}},
			},
			wantErr: false,
		},
		{
			name:    "log analysis with empty input throws error",
			entries: []LogEntry{},
			topN:    2,
			want:    nil,
			wantErr: true,
		},
		{
			name: "log analysis with topN greater than entries throws error",
			entries: []LogEntry{
				{StatusCode: 404, URL: "/home"},
				{StatusCode: 404, URL: "/about"},
				{StatusCode: 404, URL: "/home"},
			},
			topN:    3,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &CombinedLogAnalyzer{}
			got, err := l.GetLogAnalysis(tt.entries, tt.topN)
			if (err != nil) != tt.wantErr {
				t.Errorf("CombinedLogAnalyzer.GetLogAnalysis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
