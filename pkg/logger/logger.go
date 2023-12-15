package logger

import "github.com/sirupsen/logrus"

type Logger interface {
	Infof(message string, args ...interface{})
	Warningf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Fatalf(message string, args ...interface{})
}

// make a log configurable, for example
// there are 5 levels of logs, Debug, Info, Warning, Error, Fatal
// in some case you don't need to log debug/info etc.

type AppLogger struct {
	Log *logrus.Logger
}

func NewLogger() Logger {
	logger := logrus.New()
	return &AppLogger{Log: logger}
}

func (a *AppLogger) Infof(message string, args ...interface{}) {
	a.Log.Infof(message, args...)
}

func (a *AppLogger) Warningf(message string, args ...interface{}) {
	a.Log.Warnf(message, args...)
}

func (a *AppLogger) Errorf(message string, args ...interface{}) {
	a.Log.Errorf(message, args...)
}

func (a *AppLogger) Fatalf(message string, args ...interface{}) {
	a.Log.Fatalf(message, args...)
}