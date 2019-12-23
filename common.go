package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	ProgramName    = "psguard"
	ProgramTitle   = "Polls processes for resources consumption and existence"
	ProgramVersion = "v0.3.0"
)

var (
	ConfigPath = fmt.Sprintf("./%s.json", ProgramName)
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

func ReadText(path string) string {
	file, err := os.Open(path)
	if err == nil {
		content, err := ioutil.ReadAll(file)
		if err == nil {
			return string(content)
		}
	}
	return ""
}
