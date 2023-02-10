package env

import (
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var (

	// receive the timeout seconds from the environmentvariables, default is 30 sec
	TIMEOUTSECONDS float64 = func() float64 {
		if tos, err := strconv.ParseFloat(os.Getenv("TIMEOUTSECONDS"), 64); err == nil {
			return tos
		}
		return 30
	}()

	// receive the timeout seconds from the environmentvariables
	LOGLEVEL logrus.Level = func() logrus.Level {
		switch strings.ToLower(os.Getenv("LOGLEVEL")) {
		case "trace":
			return logrus.TraceLevel
		case "debug":
			return logrus.DebugLevel
		case "info":
			return logrus.InfoLevel
		case "warn":
			return logrus.WarnLevel
		case "error":
			return logrus.ErrorLevel
		case "fatal":
			return logrus.FatalLevel
		case "panic":
			return logrus.PanicLevel
		}
		return logrus.WarnLevel
	}()
)

// initialize the log provider -> logrus
var _log *logrus.Logger = func() *logrus.Logger {

	var l *logrus.Logger = &logrus.Logger{
		Out:          os.Stdout,
		Level:        LOGLEVEL,
		Formatter:    &logrus.JSONFormatter{},
		ReportCaller: true,
	}

	l.Debug("logger initialized")

	l.WithFields(logrus.Fields{
		"env.TIMEOUTSECONDS": TIMEOUTSECONDS,
		"env.LOGLEVEL":       LOGLEVEL,
	}).Info("environment variables processed")

	return l
}()

// service instance of a logrus logger
func Log() *logrus.Logger {
	return _log
}
