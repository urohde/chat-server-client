package logger

import (
	"fmt"
	"os"
)

type FileLogger struct {
	file *os.File
}

func NewFileLogger(filename string) (*FileLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[FileLogger] Logging to file: %s\n", filename)
	return &FileLogger{file}, nil
}

func (l *FileLogger) Write(msg []byte) error {
	_, err := l.file.Write(msg)
	if err != nil {
		return fmt.Errorf("error writing to log file: %w", err)
	}

	return nil
}

func (l *FileLogger) Close() error {
	return l.file.Close()
}
