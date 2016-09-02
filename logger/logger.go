package logger

import "github.com/lfq7413/tomato/types"

const logStringTruncateLength = 1000
const truncationMarker = "... (truncated)"

var adapter loggerAdapter

func init() {
	adapter = newBeegoLogger()
}

// Log ...
func Log(level string, args ...interface{}) {
	adapter.log(level, args...)
}

// Info ...
func Info(args ...interface{}) {
	Log("info", args...)
}

// Error ...
func Error(args ...interface{}) {
	Log("error", args...)
}

// Warn ...
func Warn(args ...interface{}) {
	Log("warn", args...)
}

// Verbose ...
func Verbose(args ...interface{}) {
	Log("verbose", args...)
}

// Debug ...
func Debug(args ...interface{}) {
	Log("debug", args...)
}

// Silly ...
func Silly(args ...interface{}) {
	Log("silly", args...)
}

// TruncateLogMessage ...
func TruncateLogMessage(msg string) string {
	if len(msg) > logStringTruncateLength {
		return msg[:logStringTruncateLength] + truncationMarker
	}
	return msg
}

func parseOptions(options map[string]string) types.M {
	// TODO
	return types.M{}
}

// GetLogs ...
func GetLogs(options map[string]string) (types.M, error) {
	return adapter.query(parseOptions(options))
}

type loggerAdapter interface {
	log(level string, args ...interface{})
	query(options types.M) (types.M, error)
}
