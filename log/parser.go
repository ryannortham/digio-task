package log

import (
	"fmt"
	"regexp"
	"strconv"
)

type LogEntry struct {
	IP         string
	Identity   string
	UserID     string
	Time       string
	Method     string
	URL        string
	Protocol   string
	StatusCode int
	Size       int
	Referrer   string
	UserAgent  string
}

type LogParser interface {
	ParseLogEntry(string) (LogEntry, error)
	ParseLogEntries([]string) ([]LogEntry, error)
}

type CombinedLogParser struct{}

func (p *CombinedLogParser) ParseLogEntry(line string) (LogEntry, error) {
	// Combined Log Format (CLF) regex
	const clfRgx = `^(\S+) (\S+) (\S+) \[([\w:/]+\s[+\-]\d{4})\] "(\S+) (\S+) (\S+)" (\d{3}) (\d+) "([^"]*)" "([^"]*)".*`
	clfRegex := regexp.MustCompile(clfRgx)
	logFields := clfRegex.FindStringSubmatch(line)

	// parse error
	if logFields == nil {
		return LogEntry{}, fmt.Errorf("log parsing error for line: %s", line)
	}

	// log parsed successfully
	logEntry := LogEntry{
		IP:         logFields[1],
		Identity:   logFields[2],
		UserID:     logFields[3],
		Time:       logFields[4],
		Method:     logFields[5],
		URL:        logFields[6],
		Protocol:   logFields[7],
		StatusCode: ParseInt(logFields[8]),
		Size:       ParseInt(logFields[9]),
		Referrer:   logFields[10],
		UserAgent:  logFields[11],
	}

	return logEntry, nil
}

func (p *CombinedLogParser) ParseLogEntries(logLines []string) ([]LogEntry, error) {
	var logEntries []LogEntry

	for _, line := range logLines {
		entry, err := p.ParseLogEntry(line)
		if err != nil {
			fmt.Printf("error parsing log entry, omitting: %v\n", err)
			continue
		}

		logEntries = append(logEntries, entry)
	}

	if len(logEntries) == 0 {
		return nil, fmt.Errorf("no log entries parsed successfully")
	}

	return logEntries, nil
}

func ParseInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		// try parsing as a float
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			fmt.Printf("error parsing int: %v\n", err)
			return -1
		}

		// convert float to int
		i = int(f)
	}

	return i
}
