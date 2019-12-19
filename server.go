package main

import (
	"log"
	"os"
	"time"
)

// PollServer models a poll server
type PollServer struct {
	Settings *Options
}

// NewPollServer builds a new instance of PollServer
func NewPollServer(settings Options) *PollServer {
	return &PollServer{&settings}
}

// Start creates begins server main loop
func (p PollServer) Start() {
	processInfo, err := NewProcessInfo(int32(p.Settings.Pid))
	if err != nil {
		Fail(err.Error())
	}
	var lastProc *ProcessInfo
	for {
		processInfo, err := NewProcessInfoFromExe(processInfo.Exe)
		switch {
		case err != nil:
			if p.Settings.RestartOnEnd {
				if err := lastProc.Start(); err != nil {
					Fail(err.Error())
				}
				log.Printf("Process '%s' restarted\n", lastProc.Exe)
				time.Sleep(p.Settings.Waiting)
			} else {
				os.Exit(0)
			}
		case err == nil:
			lastProc = processInfo.Clone()
			if processInfo.CPUPercent > p.Settings.CPUPercent {
				log.Printf("CPU constraint of %.2f%% violated by +%.2f%%",
					p.Settings.CPUPercent, processInfo.CPUPercent-p.Settings.CPUPercent)
			}
			if processInfo.MemoryPercent > p.Settings.MemoryPercent {
				log.Printf("Memory constraint of %.2f%% violated by +%.2f%%",
					p.Settings.MemoryPercent, processInfo.MemoryPercent-p.Settings.MemoryPercent)
			}
		}
		time.Sleep(p.Settings.Polling)
	}
}
