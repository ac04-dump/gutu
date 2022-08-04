package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexcoder04/friendly"
)

func GetDispServer() string {
	sType := os.Getenv("XDG_SESSION_TYPE")
	if sType == "" {
		if os.Getenv("DISPLAY") != "" {
			return "x11"
		}
		return ""
	}
	return sType
}

func DesktopRunning() bool {
	dispServer := GetDispServer()
	if dispServer == "" {
		return false
	}
	if dispServer == "wayland" {
		display := os.Getenv("WAYLAND_DISPLAY")
		if display == "" {
			return false
		}
		uid := os.Getuid()
		waySock := fmt.Sprintf("/run/user/%d/%s", uid, display)
		return friendly.IsFile(waySock)
	}
	display := strings.Replace(os.Getenv("DISPLAY"), ":", "", 1)
	if display == "" {
		return false
	}
	x11Sock := fmt.Sprintf("/tmp/.X11-unix/X%s", display)
	return friendly.IsFile(x11Sock)
}
