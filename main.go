package main

import (
	"flag"
	"os"
	"sync"

	"github.com/alexcoder04/friendly/linux"
	"github.com/sirupsen/logrus"
)

const PROGRAM_NAME = "gutu"

var (
	ConfigFolder = flag.String("config", GetConfigFolder(), "config folder")
	LogFolder    = flag.String("log", GetLogFolder(), "log folder")
	Logger       *logrus.Logger
)

func main() {
	flag.Parse()

	InitLogger()

	if !linux.GuiRunning() {
		Logger.WithField("service", "_main").Fatal("Desktop not running")
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
