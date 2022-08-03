package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

func GetConfigDir() string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln("[FATL] Cannot get user config dir")
	}
	configDir := path.Join(userConfigDir, PROGRAM_NAME)
	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		log.Fatalln("[FATL] Cannot create config dir")
	}
	return configDir
}

func GetLogDir() string {
	logDir := path.Join(os.TempDir(), PROGRAM_NAME)
	err := os.MkdirAll(logDir, 0700)
	if err != nil {
		log.Fatalln("[FATL] Cannot create log dir")
	}
	return logDir
}

func GetLogFileForService(name string) string {
	return path.Join(GetLogDir(), fmt.Sprintf("%s-%s.log", name, time.Now().Format("060102-150405")))
}

func GetLogFile() string {
	return path.Join(GetLogDir(), fmt.Sprintf("%s.log", time.Now().Format("060102-150405")))
}
