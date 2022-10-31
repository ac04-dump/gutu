package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/alexcoder04/friendly/v2/flinux"
	"github.com/mitchellh/go-ps"
	"gopkg.in/yaml.v3"
)

type Service struct {
	Name        string   `yaml:"Name"`
	Command     string   `yaml:"Command"`
	Args        []string `yaml:"Args"`
	Interval    uint     `yaml:"Interval"`
	KeepAlive   bool     `yaml:"KeepAlive"`
	RetryNumber uint     `yaml:"RetryNumber"`
	When        string   `yaml:"When"`
	Delay       uint     `yaml:"Delay"`
	KillOld     bool     `yaml:"KillOld"`
}

func GetServices() []Service {
	files, err := ioutil.ReadDir(*ConfigFolder)
	if err != nil {
		Logger.WithField("service", "_main").Fatalln("Cannot read contents of config dir")
	}

	services := []Service{}
	for _, f := range files {
		// ignore imcompatible extensions
		extension := path.Ext(f.Name())
		if extension != ".service" && extension != ".yml" && extension != ".yaml" {
			continue
		}

		// read file
		data, err := ioutil.ReadFile(path.Join(*ConfigFolder, f.Name()))
		if err != nil {
			Logger.WithField("service", "_main").Warnf("Cannot read %s, skipping\n", f.Name())
			continue
		}

		// unmarshal data
		s := Service{}
		err = yaml.Unmarshal(data, &s)
		if err != nil {
			Logger.WithField("service", f.Name()).Warn("Cannot unmarshal, skipping")
			continue
		}

		if s.When == "never" {
			Logger.WithField("service", s.Name).Info("Disabled, skipping")
			continue
		}
		if s.When != flinux.GetDisplayServer() && s.When != "always" {
			Logger.WithField("service", s.Name).Info("Disabled on this display server, skipping")
			continue
		}

		services = append(services, s)
		Logger.WithField("service", s.Name).Info("Loaded")
	}
	return services
}

func GetPidByName(name string) int {
	processes, err := ps.Processes()
	if err != nil {
		return 0
	}
	for _, p := range processes {
		if p.Executable() == name {
			return p.Pid()
		}
	}
	return 0
}

func HandleService(s Service, wg *sync.WaitGroup) {
	defer wg.Done()

	if s.Delay > 0 {
		Logger.WithField("service", s.Name).Infof("Delaying by %ds", s.Delay)
		time.Sleep(time.Duration(s.Delay) * time.Second)
	}

	f, err := os.OpenFile(GetLogFileForService(s.Name), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		Logger.WithField("service", s.Name).Warn("Failed to open log file, skipping")
		return
	}

	if s.KillOld {
		pid := GetPidByName(s.Command)
		if pid != 0 {
			proc, err := os.FindProcess(pid)
			if err == nil {
				proc.Kill()
			}
		}
	}

	for i := 0; true; i++ {
		if i > int(s.RetryNumber) && s.Interval == 0 {
			Logger.WithField("service", s.Name).Warn("Retry number exceeded, stopping")
			return
		}

		var stdBuf bytes.Buffer
		mw := io.MultiWriter(f, &stdBuf)
		command := exec.Command(s.Command, s.Args...)
		command.Stdout = mw
		command.Stderr = mw

		Logger.WithField("service", s.Name).Info("Starting")
		err := command.Start()
		if err != nil {
			Logger.WithField("service", s.Name).Warn("Failed to start, retrying in %ds\n", s.RetryNumber)
			time.Sleep(time.Duration(s.Interval) * time.Second)
		}
		err = command.Wait()
		if err != nil {
			Logger.WithField("service", s.Name).Warnf("Error waiting for command to finish (%s), restarting in %ds", err.Error(), s.RetryNumber)
			time.Sleep(time.Duration(s.Interval) * time.Second)
		}

		if !flinux.GuiRunning() {
			Logger.WithField("service", s.Name).Warn("Desktop seems to be shut down, stopping")
			return
		}

		if !s.KeepAlive && s.Interval == 0 {
			Logger.WithField("service", s.Name).Info("Finished")
			return
		}

		Logger.WithField("service", s.Name).Infof("Restarting in %ds", s.Interval)
		time.Sleep(time.Duration(s.Interval) * time.Second)
	}
}
