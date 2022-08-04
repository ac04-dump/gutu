package main

import (
	"fmt"
	"os"

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
			display = "1"
		}
		uid := os.Getuid()
		waySock := fmt.Sprintf("/run/user/%d/wayland-%s", uid, display)
		return friendly.IsFile(waySock)
	}
	display := os.Getenv("DISPLAY")
	if display == "" {
		display = "0"
	}
	x11Sock := fmt.Sprintf("/tmp/.X11-unix/X%s", display)
	return friendly.IsFile(x11Sock)
}
