package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var LevelMap = map[string]string{
	"trace":   "TRACE",
	"debug":   "DEBUG",
	"info":    "INFO",
	"warning": "WARN",
	"error":   "ERR",
	"fatal":   "FATAL",
	"panic":   "PANIC",
}

type GutuLogFormatter struct {
}

func (f *GutuLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	str := fmt.Sprintf(
		"[%s] [%s] [%s] %s\n",
		entry.Time.Format("2006-01-02-15:04:05"),
		LevelMap[entry.Level.String()],
		entry.Data["service"],
		entry.Message)
	return []byte(str), nil
}

func GetLogWriter() io.Writer {
	f, err := os.OpenFile(GetLogFile(), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic("Cannot open log file")
	}
	return f
}

func InitLogger() {
	Logger = &logrus.Logger{
		Out:       GetLogWriter(),
		Level:     logrus.DebugLevel,
		Formatter: &GutuLogFormatter{},
	}
}
