package livequery

import "fmt"

var logLevel = map[string]int{
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
		level: "INFO",
	}
}

type tLog struct {
	level string
}

func (l *tLog) getCurrentLogLevel() int {
	if level, ok := logLevel[l.level]; ok {
		return level
	}
	return logLevel["ERROR"]
}

func (l *tLog) verbose(args ...interface{}) {
	if l.getCurrentLogLevel() <= logLevel["VERBOSE"] {
		fmt.Println(args...)
	}
}

func (l *tLog) log(args ...interface{}) {
	if l.getCurrentLogLevel() <= logLevel["INFO"] {
		fmt.Println(args...)
	}
}

func (l *tLog) error(args ...interface{}) {
	if l.getCurrentLogLevel() <= logLevel["ERROR"] {
		fmt.Println(args...)
	}
}
