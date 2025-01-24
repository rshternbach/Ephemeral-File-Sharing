package utils

import (
	"os"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetLevel(logrus.InfoLevel)

	// Create logs folder if it doesn't exist
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}

	// Configure log file
	logFile, err := os.OpenFile("logs/file-sharing-service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Fatalf("Failed to open log file: %v", err)
	}

	// Configure logger to write to file and console
	Log.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  logFile,
			logrus.WarnLevel:  logFile,
			logrus.ErrorLevel: logFile,
			logrus.FatalLevel: logFile,
			logrus.PanicLevel: logFile,
		},
		&logrus.TextFormatter{
			FullTimestamp: true,
		},
	))
}
