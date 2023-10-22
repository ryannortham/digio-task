package log

import (
	"bufio"
	"os"
)

type LogReader interface {
	ReadLines() ([]string, error)
}

type FileReader struct {
	Filename string
}

// ReadLines reads a whole file into memory as a slice of strings.
func (r FileReader) ReadLines() ([]string, error) {
	file, err := os.Open(r.Filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
