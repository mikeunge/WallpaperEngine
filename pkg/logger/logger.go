package logger

import (
	"fmt"
	"os"

	"github.com/mikeunge/WallpaperEngine/pkg/helpers"
	"github.com/sirupsen/logrus"
)

var (
	log         *logrus.Logger
	logFilePath = "~/.local/share/wallpaper-engine.log"
)

func SetLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.WarnLevel)
	}
}

func SetOutput(fileOut bool) {
    logFilePath = helpers.SanitizePath(logFilePath)
    if fileOut {
        file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            panic(fmt.Sprintf("error opening file: %v", err))
        }
        log.SetOutput(file)
    } else {
        log.SetOutput(os.Stdout)
    }
}

func init() {
	// set up a new json logger, with loglevel warn and write to file
	log = logrus.New()
	log.Formatter = &logrus.JSONFormatter{}
	SetLogLevel("warn")
    SetOutput(true)
}

func Debug(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func Info(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func Warn(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

func Error(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

var (
	ConfigError = "%v type=config.error"
)
