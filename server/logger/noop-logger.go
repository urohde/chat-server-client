package logger

import "fmt"

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	fmt.Println("[NoopLogger] Logging disabled")
	return &NoopLogger{}
}

func (l *NoopLogger) Write(msg []byte) error {
	return nil
}

func (l *NoopLogger) Close() error {
	return nil
}
