package main

import (
	"os"
)

func main() {
	options := NewOptions()
	if options == nil {
		PrintInfo()
		os.Exit(0)
	}
	if err := options.Validate(); err != nil {
		Exit(err.Error(), true)
	}

	server := NewPollServer(*options)
	server.Start()
}
