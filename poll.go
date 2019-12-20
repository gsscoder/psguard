package main

import (
	"log"
	"time"
)

// PollManager models a poll server
type PollManager struct {
	Settings   *Options
	StopSignal chan struct{}
}

// NewPollManager builds a new instance of PollServer
func NewPollManager(settings Options) *PollManager {
	return &PollManager{Settings: &settings}
}

// Start begins the polling loop
func (p PollManager) Start() {
	p.StopSignal = make(chan struct{})

	process, err := NewProcessInfo(int32(p.Settings.Pid))
	if err != nil {
		Fail(err.Error())
	}
	go p.poll(*process, *p.Settings)
	<-p.StopSignal
	close(p.StopSignal)
}

func (p PollManager) poll(process ProcessInfo, settings Options) {
	var lastProcess *ProcessInfo
	for {
		process, err := NewProcessInfoFromExe(process.Exe)
		switch err {
		case nil:
			lastProcess = process.Clone()
			if process.CPUPercent > p.Settings.CPUPercent {
				log.Printf("CPU constraint of %.2f%% violated by +%.2f%%",
					p.Settings.CPUPercent, process.CPUPercent-p.Settings.CPUPercent)
			}
			if process.MemoryPercent > p.Settings.MemoryPercent {
				log.Printf("Memory constraint of %.2f%% violated by +%.2f%%",
					p.Settings.MemoryPercent, process.MemoryPercent-p.Settings.MemoryPercent)
			}
		default:
			if p.Settings.RestartOnEnd {
				if err := lastProcess.Start(); err != nil {
					log.Fatal(err.Error())
					p.stop()
				}
				log.Printf("Process '%s' restarted\n", lastProcess.Exe)
				time.Sleep(p.Settings.Waiting)
			} else {
				p.stop()
			}
		}
		time.Sleep(p.Settings.Polling)
	}
}

func (p PollManager) stop() {
	p.StopSignal <- struct{}{}
}
