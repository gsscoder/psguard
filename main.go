package main

import (
	"log"
	"os"
)

func main() {
	options := NewOptions()
	if options == nil {
		PrintInfo()
		os.Exit(0)
	}
	config := ReadText(ConfigPath)
	if len(config) == 0 {
		log.Fatal("Unable to read config file")
	}
	constraints := NewConstraintList(config)
	if len(constraints) == 0 {
		log.Fatal("Process list is empty")
	}
	if !constraints.HasProcesses() {
		log.Fatal("Unable to find any process")
	}
	messagges := constraints.Sanitize(*options)
	for _, message := range messagges {
		log.Print(message)
	}
	server := NewPollManager(*options, constraints)
	server.Start()
}
