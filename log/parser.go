package log

import (
	"fmt"
	"regexp"
	"strconv"
)

type LogParser interface {
	ParseLogEntry(string) (LogEntry, error)
	ParseLogEntries([]string) ([]LogEntry, error)
}

type CombinedLogParser struct{}

// ParseLogEntry parses a log entry string into a LogEntry struct
func (p *CombinedLogParser) ParseLogEntry(line string) (LogEntry, error) {
	// Common Log Format (CLF) regex
	const clfRgx = `^(\S+) (\S+) (\S+) \[([\w:/]+\s[+\-]\d{4})\] "(\S+) (\S+) (\S+)" (\d{3}) (\d+) "([^"]*)" "([^"]*)".*`
	clfRegex := regexp.MustCompile(clfRgx)

	var logFields []string
	if logFields = clfRegex.FindStringSubmatch(line); logFields == nil {
		return LogEntry{}, fmt.Errorf("log parsing error for line: %s", line)
	}

	entry := LogEntry{
		IP:         logFields[1],
		Identity:   logFields[2],
		UserID:     logFields[3],
		Time:       logFields[4],
		Method:     logFields[5],
		Path:       logFields[6],
		Protocol:   logFields[7],
		StatusCode: parseInt(logFields[8]),
		Size:       parseInt(logFields[9]),
		Referrer:   logFields[10],
		UserAgent:  logFields[11],
	}

	return entry, nil
}

// ParseLogEntries parses a slice of log entry strings into a slice of LogEntry structs
func (p *CombinedLogParser) ParseLogEntries(logLines []string) ([]LogEntry, error) {
	var logEntries []LogEntry
	for _, line := range logLines {
		entry, err := p.ParseLogEntry(line)
		if err != nil {
			fmt.Printf("Error parsing log entry, omitting: %v\n", err)
			continue
		}

		logEntries = append(logEntries, entry)
	}

	return logEntries, nil
}

// parseInt parses a string into an int
func parseInt(intStr string) int {
	i, err := strconv.Atoi(intStr)
	if err != nil {
		fmt.Printf("error parsing int: %v\n", err)
		return -1
	}

	return i
}
