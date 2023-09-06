package sorter

// Logger that is able to send log messages from the sorter instance
type SorterLogger interface {
	// Log a message with a given format with the Debug log level
	Debugf(format string, args ...interface{})

	// Log a message with a given format with the Info log level
	Infof(format string, args ...interface{})

	// Log a message with a given format with the Warning log level
	Warnf(format string, args ...interface{})

	// Log a message with a given format with the Error log level
	Errorf(format string, args ...interface{})
}

type discardLogger struct{}

func (dl *discardLogger) Debugf(_ string, _ ...interface{}) {}
func (dl *discardLogger) Infof(_ string, _ ...interface{})  {}
func (dl *discardLogger) Warnf(_ string, _ ...interface{})  {}
func (dl *discardLogger) Errorf(_ string, _ ...interface{}) {}

// Create a new sorter logger instance. The logs are not written anywhere. This logger can be used
// for mocking purposes or as a subsitute for a provided nil reference logger
func getDiscardLogger() SorterLogger {
	return new(discardLogger)
}
