package main

import (
	"flag"
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

const PROGRAM_NAME = "gutu"

var (
	ConfigFolder = flag.String("config", GetConfigDir(), "config folder")
	LogFolder    = flag.String("log", GetLogDir(), "log folder")
	Logger       *logrus.Logger
)

func GetLogWriter() io.Writer {
	f, err := os.OpenFile(GetLogFile(), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic("Cannot open log file")
	}
	return f
}

func main() {
	flag.Parse()
	Logger = &logrus.Logger{
		Out:   GetLogWriter(),
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
	}

	oldPid := GetPidByName("gutu")
	if oldPid != 0 {
		proc, err := os.FindProcess(oldPid)
		if err == nil && oldPid != os.Getpid() {
			Logger.WithField("service", "_main").Info("Killing old instance")
			proc.Kill()
		}
	}

	Logger.WithField("service", "_main").Info("Loading and starting services")
	var wg sync.WaitGroup

	for _, s := range GetServices() {
		wg.Add(1)
		go HandleService(s, &wg)
	}

	wg.Wait()
	Logger.WithField("service", "_main").Info("Done")
}
