package main

import (
	"os"
)

func main() {
	opts := NewOptions()
	if opts == nil {
		programInfo()
		os.Exit(0)
	}
	err := opts.Validate()
	if err != nil {
		Exit(err.Error(), true)
	}
	NewApp(opts)
}
