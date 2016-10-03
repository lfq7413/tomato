package utils

import "fmt"

// LogLevel ...
var LogLevel = map[string]int{
	"VERBOSE": 0,
	"DEBUG":   1,
	"INFO":    2,
	"ERROR":   3,
	"NONE":    4,
}

// TLog ...
var TLog *tLog

func init() {
	TLog = &tLog{
		Level: "INFO",
	}
}

type tLog struct {
	Level string
}

func (l *tLog) getCurrentLogLevel() int {
	if level, ok := LogLevel[l.Level]; ok {
		return level
	}
	return LogLevel["ERROR"]
}

func (l *tLog) Verbose(args ...interface{}) {
	if l.getCurrentLogLevel() <= LogLevel["VERBOSE"] {
		fmt.Println(args...)
	}
}

func (l *tLog) Log(args ...interface{}) {
	if l.getCurrentLogLevel() <= LogLevel["INFO"] {
		fmt.Println(args...)
	}
}

func (l *tLog) Error(args ...interface{}) {
	if l.getCurrentLogLevel() <= LogLevel["ERROR"] {
		fmt.Println(args...)
	}
}
