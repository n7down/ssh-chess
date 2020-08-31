package logruslogger

import (
	log "github.com/sirupsen/logrus"
)

type LogrusLogger struct{}

func NewLogrusLogger(showMethod bool) *LogrusLogger {
	log.SetReportCaller(showMethod)
	return &LogrusLogger{}
}

func (l LogrusLogger) Trace(args ...interface{}) {
	log.Trace(args...)
}

func (l LogrusLogger) Debug(args ...interface{}) {
	log.Debug(args...)
}

func (l LogrusLogger) Print(args ...interface{}) {
	log.Print(args...)
}

func (l LogrusLogger) Info(args ...interface{}) {
	log.Info(args...)
}

func (l LogrusLogger) Warn(args ...interface{}) {
	log.Warn(args...)
}

func (l LogrusLogger) Warning(args ...interface{}) {
	log.Warning(args...)
}

func (l LogrusLogger) Error(args ...interface{}) {
	log.Error(args...)
}

func (l LogrusLogger) Panic(args ...interface{}) {
	log.Panic(args...)
}

func (l LogrusLogger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func (l LogrusLogger) Tracef(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

func (l LogrusLogger) Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func (l LogrusLogger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l LogrusLogger) Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func (l LogrusLogger) Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func (l LogrusLogger) Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

func (l LogrusLogger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func (l LogrusLogger) Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func (l LogrusLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func (l LogrusLogger) Traceln(args ...interface{}) {
	log.Traceln(args...)
}

func (l LogrusLogger) Debugln(args ...interface{}) {
	log.Debugln(args...)
}

func (l LogrusLogger) Println(args ...interface{}) {
	log.Println(args...)
}

func (l LogrusLogger) Infoln(args ...interface{}) {
	log.Infoln(args...)
}

func (l LogrusLogger) Warnln(args ...interface{}) {
	log.Warnln(args...)
}

func (l LogrusLogger) Warningln(args ...interface{}) {
	log.Warningln(args...)
}

func (l LogrusLogger) Errorln(args ...interface{}) {
	log.Errorln(args...)
}

func (l LogrusLogger) Panicln(args ...interface{}) {
	log.Panicln(args...)
}

func (l LogrusLogger) Fatalln(args ...interface{}) {
	log.Fatalln(args...)
}
