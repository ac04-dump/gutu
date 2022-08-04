package main

import (
	"fmt"
	"os"
	"path"
	"time"
)

// runs before main
func GetConfigDir() string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		Logger.WithField("service", "_util").Fatalln("[FATL] Cannot get user config dir")
	}
	configDir := path.Join(userConfigDir, PROGRAM_NAME)
	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		Logger.WithField("service", "_util").Fatalln("[FATL] Cannot create config dir")
	}
	return configDir
}

// runs before main
func GetLogDir() string {
	logDir := path.Join(os.TempDir(), fmt.Sprintf("%s-%d", PROGRAM_NAME, os.Getuid()))
	err := os.MkdirAll(logDir, 0700)
	if err != nil {
		Logger.WithField("service", "_util").Fatalln("[FATL] Cannot create log dir")
	}
	return logDir
}

func GetLogFileForService(name string) string {
	return path.Join(*LogFolder, fmt.Sprintf("%s-%s.log", name, time.Now().Format("060102-150405")))
}

func GetLogFile() string {
	return path.Join(*LogFolder, fmt.Sprintf("%s.log", time.Now().Format("060102-150405")))
}
