package main

import (
	"log"
	"os"
	"sync"
)

const PROGRAM_NAME = "gutu"

func main() {
	f, err := os.OpenFile(GetLogFile(), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic("Cannot open log file")
	}
	defer f.Close()

	log.SetOutput(f)

	log.Println("[INFO] [_main] Loading and starting services")
	var wg sync.WaitGroup

	for _, s := range GetServices() {
		wg.Add(1)
		go HandleService(s, &wg)
	}

	wg.Wait()
	log.Println("[INFO] [_main] Done")
}
