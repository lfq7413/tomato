package logger

import (
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
)

type beegoLogger struct {
	beelogger *logs.BeeLogger
}

func newBeegoLogger() *beegoLogger {
	l := logs.NewLogger(1000)
	l.SetLevel(logs.LevelDebug)
	l.SetLogger("file", `{"filename":"project.log"}`)
	l.DelLogger("console")
	l.Async()
	return &beegoLogger{
		beelogger: l,
	}
}

func (l *beegoLogger) log(level string, args ...interface{}) {
	switch level {
	case "debug":
		l.beelogger.Debug(generateFmtStr(len(args)), args...)

	case "info":
		l.beelogger.Informational(generateFmtStr(len(args)), args...)

	case "verbose":
		l.beelogger.Notice(generateFmtStr(len(args)), args...)

	case "warn":
		l.beelogger.Warning(generateFmtStr(len(args)), args...)

	case "error":
		l.beelogger.Error(generateFmtStr(len(args)), args...)

	case "silly":
		l.beelogger.Critical(generateFmtStr(len(args)), args...)

	}
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}

func (l *beegoLogger) query(options types.M) (types.M, error) {
	return nil, errs.E(errs.PushMisconfigured, "Querying logs is not supported with this adapter")
}
