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
	constraints := NewConstraintList()
	if len(constraints) == 0 {
		Fail("Process list is empty")
	}
	if !constraints.HasProcesses() {
		Fail("Unable to bind any process")
	}
	messagges := constraints.Sanitize(*options)
	for _, message := range messagges {
		log.Print(message)
	}
	server := NewPollManager(*options, constraints)
	server.Start()
}
