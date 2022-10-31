package main

import (
	"fmt"
	"path"
	"time"

	"github.com/alexcoder04/friendly/v2/ffiles"
)

// get executed before main() {{{
func GetConfigFolder() string {
	dir, err := ffiles.GetConfigDirFor(PROGRAM_NAME)
	if err != nil {
		panic("Cannot get config folder")
	}
	return dir
}

func GetLogFolder() string {
	dir, err := ffiles.GetLogDirFor(PROGRAM_NAME)
	if err != nil {
		panic("Cannot get log folder")
	}
	return dir
}

// }}}

func GetLogFileForService(name string) string {
	return path.Join(*LogFolder, fmt.Sprintf("%s-%s.log", name, time.Now().Format("060102-150405")))
}

func GetLogFile() string {
	return path.Join(*LogFolder, fmt.Sprintf("%s.log", time.Now().Format("060102-150405")))
}
