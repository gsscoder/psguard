package main

import (
	"log"
	"time"
)

// PollManager models a poll server
type PollManager struct {
	Polling      time.Duration
	RestartOnEnd bool
	Waiting      time.Duration
	Constraints  []Constraint
	// StopSignal   chan struct{}
}

// NewPollManager builds a new instance of PollServer
func NewPollManager(settings Options, constr []Constraint) *PollManager {
	return &PollManager{
		Polling:      settings.Polling,
		RestartOnEnd: settings.RestartOnEnd,
		Waiting:      settings.Waiting,
		Constraints:  constr}
}

// Start begins the polling loop
func (p PollManager) Start() {
	// p.StopSignal = make(chan struct{})
	for {
		p.poll()
		time.Sleep(p.Polling)
	}
	//<-p.StopSignal
}

func (p PollManager) poll() {
	for i := range p.Constraints {
		go p.pollProcess(&p.Constraints[i])
	}
}

func (p PollManager) pollProcess(constr *Constraint) {
	if constr.Process == nil {
		proc, err := NewProcessInfo(constr.Pid)
		if err == nil {
			constr.Process = proc
		} else {
			log.Printf("Can't find process %s", constr.Name)
			return
		}
	}
	proc, err := NewProcessInfoFromExe(constr.Process.Exe)
	switch err {
	case nil:
		if proc.CPUPercent > constr.CPUPercent {
			log.Printf("%s: CPU constraint of %.2f%% violated by +%.2f%%",
				constr.Name, constr.CPUPercent, proc.CPUPercent-constr.CPUPercent)
		}
		if proc.MemoryPercent > constr.MemoryPercent {
			log.Printf("%s: Memory constraint of %.2f%% violated by +%.2f%%",
				constr.Name, constr.MemoryPercent, proc.MemoryPercent-constr.MemoryPercent)
		}
	default:
		log.Printf("Process %s terminated", constr.Name)
		if p.RestartOnEnd {
			if err := constr.Process.Start(); err != nil {
				log.Fatal(err.Error())
			}
			log.Printf("Process %s restarted\n", constr.Process.Exe)
			time.Sleep(p.Waiting)
		}
	}
}
