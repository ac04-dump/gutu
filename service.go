package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

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
	configDir := GetConfigDir()
	files, err := ioutil.ReadDir(configDir)
	if err != nil {
		log.Fatalln("[FATL] Cannot read contents of config dir")
	}

	services := []Service{}
	for _, f := range files {
		// ignore imcompatible extensions
		extension := path.Ext(f.Name())
		if extension != ".service" && extension != ".yml" && extension != ".yaml" {
			continue
		}

		// read file
		data, err := ioutil.ReadFile(path.Join(configDir, f.Name()))
		if err != nil {
			log.Printf("[WARN] Cannot read %s, skipping\n", f.Name())
			continue
		}

		// unmarshal data
		s := Service{}
		err = yaml.Unmarshal(data, &s)
		if err != nil {
			log.Printf("[WARN] [%s] Cannot unmarshal, skipping\n", f.Name())
			continue
		}

		if s.When == "never" {
			log.Printf("[INFO] [%s] Disabled, skipping\n", s.Name)
			continue
		}
		if s.When != GetDispServer() && s.When != "always" {
			log.Printf("[INFO] [%s] Disabled on this display server, skipping\n", s.Name)
			continue
		}

		services = append(services, s)
		log.Printf("[INFO] [%s] Loaded", s.Name)
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
		log.Printf("[INFO] [%s] Delaying by %ds\n", s.Name, s.Delay)
		time.Sleep(time.Duration(s.Delay) * time.Second)
	}

	f, err := os.OpenFile(GetLogFileForService(s.Name), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("[WARN] [%s] Failed to open log file, skipping", s.Name)
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
		if i > int(s.RetryNumber) {
			log.Printf("[WARN] [%s] Retry number exceeded, stopping\n", s.Name)
			return
		}

		var stdBuf bytes.Buffer
		mw := io.MultiWriter(f, &stdBuf)
		command := exec.Command(s.Command, s.Args...)
		command.Stdout = mw
		command.Stderr = mw

		err := command.Start()
		if err != nil {
			log.Printf("[WARN] [%s] Failed to start, retrying in %ds\n", s.Name, s.RetryNumber)
			time.Sleep(time.Duration(s.Interval) * time.Second)
		}
		err = command.Wait()
		if err != nil {
			log.Printf("[WARN] [%s] Error waiting for command to finish, restarting in %ds\n", s.Name, s.RetryNumber)
			time.Sleep(time.Duration(s.Interval) * time.Second)
		}

		if !DesktopRunning() {
			log.Printf("[WARN] [%s] Desktop seems to be shut down, stopping\n", s.Name)
			return
		}

		if !s.KeepAlive && s.Interval == 0 {
			log.Printf("[INFO] [%s] Finished\n", s.Name)
			return
		}

		log.Printf("[INFO] [%s] Restarting in %ds", s.Name, s.Interval)
		time.Sleep(time.Duration(s.Interval) * time.Second)
	}
}
