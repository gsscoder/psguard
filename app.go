package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	programName      = "psguard"
	programTitle     = "Polls a process for resources consumption and existence"
	programVersion   = "v0.1.0"
	defCPUPercent    = 1
	defMemoryPercent = 1
	defPolling       = time.Second
	defWaiting       = time.Second * 5
)

// NewApp creates a new application with given command line options
func NewApp(opts *Options) {
	processInfo, err := NewProcessInfo(int32(opts.Pid))
	if err != nil {
		fail(err.Error())
	}
	var lastProc *ProcessInfo
	for {
		processInfo, err := NewProcessInfoFromExe(processInfo.Exe)
		switch {
		case err != nil:
			if opts.RestartOnEnd {
				if err := lastProc.Start(); err != nil {
					fail(err.Error())
				}
				log.Printf("Process '%s' restarted\n", lastProc.Exe)
				time.Sleep(opts.Waiting)
			} else {
				os.Exit(0)
			}
		case err == nil:
			lastProc = processInfo.Clone()
			if processInfo.CPUPercent > opts.CPUPercent {
				log.Printf("CPU constraint of %.2f%% violated by +%.2f%%",
					opts.CPUPercent, processInfo.CPUPercent-opts.CPUPercent)
			}
			if processInfo.MemoryPercent > opts.MemoryPercent {
				log.Printf("Memory constraint of %.2f%% violated by +%.2f%%",
					opts.MemoryPercent, processInfo.MemoryPercent-opts.MemoryPercent)
			}
		}
		time.Sleep(opts.Polling)
	}
}

// Exit terminates with a message and optionally suggestion to activate help
func Exit(message string, help bool) {
	fmt.Printf("%s: %s\n", programName, message)
	if help {
		fmt.Printf("Try %s --help\n", programName)
	}
	os.Exit(1)
}

func programInfo() {
	fmt.Printf("%s: %s\n", programName, programTitle)
	fmt.Printf("Version: %s\n", programVersion)
}

func fail(message string) {
	log.Fatal(message)
	os.Exit(1)
}
