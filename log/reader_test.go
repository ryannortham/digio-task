package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileReader_ReadLines(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		want        []string
		wantErr     bool
	}{
		{
			name:        "read valid file",
			fileContent: "line 1\nline 2\nline 3\n",
			want:        []string{"line 1", "line 2", "line 3"},
			wantErr:     false,
		},
		{
			name:        "read empty file",
			fileContent: "",
			want:        []string{},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file with the test data
			tmpfile, err := os.CreateTemp("", "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name()) // clean up

			if _, err := tmpfile.Write([]byte(tt.fileContent)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// Replace the log file path with the temporary file path
			reader := &FileReader{LogFilePath: tmpfile.Name()}
			got, err := reader.ReadLines()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
