package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CombinedLogParser_ParseLogEntry(t *testing.T) {
	parser := &CombinedLogParser{}

	tests := []struct {
		name    string
		line    string
		want    LogEntry
		wantErr bool
	}{
		{
			name: "parse valid log entry",
			line: `127.0.0.1 - - [01/Jan/2022:00:00:00 +0000] "GET / HTTP/1.1" 200 1234 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"`,
			want: LogEntry{
				IP:         "127.0.0.1",
				Identity:   "-",
				UserID:     "-",
				Time:       "01/Jan/2022:00:00:00 +0000",
				Method:     "GET",
				URL:        "/",
				Protocol:   "HTTP/1.1",
				StatusCode: 200,
				Size:       1234,
				Referrer:   "-",
				UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
			},
			wantErr: false,
		},
		{
			name:    "parse invalid log entry throws error",
			line:    "invalid log entry",
			want:    LogEntry{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseLogEntry(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("CombinedLogParser.ParseLogEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_ParseInt(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    int
		wantErr bool
	}{
		{
			name:    "parse integer",
			str:     "123",
			want:    123,
			wantErr: false,
		},
		{
			name:    "parse float",
			str:     "123.456",
			want:    123,
			wantErr: false,
		},
		{
			name:    "parse invalid string",
			str:     "invalid",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInt(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_CombinedLogParser_ParseLogEntries(t *testing.T) {
	parser := &CombinedLogParser{}

	tests := []struct {
		name     string
		logLines []string
		wantLen  int
		wantErr  bool
	}{
		{
			name: "parse valid log lines",
			logLines: []string{
				"127.0.0.1 - - [01/Jan/2022:00:00:00 +0000] \"GET / HTTP/1.1\" 200 1234 \"-\" \"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3\"",
				"127.0.0.1 - - [01/Jan/2022:00:00:01 +0000] \"GET /about HTTP/1.1\" 200 5678 \"-\" \"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3\"",
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "parse invalid log lines",
			logLines: []string{
				"invalid log line",
				"another invalid log line",
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseLogEntries(tt.logLines)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLogEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Len(t, got, tt.wantLen)
		})
	}
}
