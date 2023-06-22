package logger

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func Debug(msg ...interface{}) {
	logrus.Debug(msg...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Info(msg ...interface{}) {
	logrus.Info(msg...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warn(msg ...interface{}) {
	logrus.Warn(msg...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Error(msg ...interface{}) {
	logrus.Error(msg...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatal(msg ...interface{}) {
	logrus.Fatal(msg...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func PrettyStruct(msg ...interface{}) {
	json, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	logrus.Info(string(json))
}
