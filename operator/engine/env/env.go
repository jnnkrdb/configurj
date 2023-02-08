package env

import (
	"os"
	"strconv"

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
		switch os.Getenv("LOGLEVEL") {
		case "Trace":
			return logrus.TraceLevel
		case "Debug":
			return logrus.DebugLevel
		case "Info":
			return logrus.InfoLevel
		case "Warn":
			return logrus.WarnLevel
		case "Error":
			return logrus.ErrorLevel
		case "Fatal":
			return logrus.FatalLevel
		case "Panic":
			return logrus.PanicLevel
		}
		return logrus.WarnLevel
	}()
)

// initialize the log provider -> logrus
var _log *logrus.Logger = func() *logrus.Logger {

	var l *logrus.Logger = &logrus.Logger{}

	l.SetLevel(LOGLEVEL)

	return l
}()

// service instance of a logrus logger
func Log() *logrus.Logger {
	return _log
}
