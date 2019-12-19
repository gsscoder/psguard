package main

import (
	"fmt"
	"log"
	"os"
)

const (
	ProgramName    = "psguard"
	ProgramTitle   = "Polls a process for resources consumption and existence"
	ProgramVersion = "v0.1.0"
)

func Exit(message string, help bool) {
	fmt.Printf("%s: %s\n", ProgramName, message)
	if help {
		fmt.Printf("Try %s --help\n", ProgramName)
	}
	os.Exit(1)
}

func PrintInfo() {
	fmt.Printf("%s: %s\n", ProgramName, ProgramTitle)
	fmt.Printf("Version: %s\n", ProgramVersion)
}

func Fail(message string) {
	log.Fatal(message)
	os.Exit(1)
}
