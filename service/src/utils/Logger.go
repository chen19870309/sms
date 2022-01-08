package utils

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  "./info.log",
		logrus.ErrorLevel: "./info.log",
	}
	Log = logrus.New()
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	return Log
}

func init() {
	Log = NewLogger()
	// Log as JSON instead of the default ASCII formatter.
	//Log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//Log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//Log.SetLevel(log.WarnLevel)
}
