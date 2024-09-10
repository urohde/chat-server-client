package logger

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (l *NoopLogger) Write(msg []byte) error {
	return nil
}

func (l *NoopLogger) Close() error {
	return nil
}
