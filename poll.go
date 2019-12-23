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
	Terminated   chan struct{}
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
	p.Terminated = make(chan struct{})

polling:
	go p.poll()
	<-p.Terminated
	if p.allTerminated() {
		close(p.Terminated)
	} else {
		goto polling
	}
}

func (p PollManager) poll() {
	for {
		for _, constr := range p.Constraints {
			p.pollProcess(constr)
		}
		time.Sleep(p.Polling)
	}
}

func (p PollManager) pollProcess(constr Constraint) {
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
		constr.Process.Terminated = true
		log.Printf("Process %s terminated", constr.Name)
		if p.RestartOnEnd {
			if err := constr.Process.Start(); err != nil {
				log.Fatal(err.Error())
				p.stop()
			}
			constr.Process.Terminated = false
			log.Printf("Process %s restarted\n", constr.Process.Exe)
			time.Sleep(p.Waiting)
		} else {
			p.stop()
		}
	}
}

func (p PollManager) allTerminated() bool {
	n := 0
	for _, constr := range p.Constraints {
		if constr.Process.Terminated {
			n++
		}
	}
	return n == len(p.Constraints)
}

func (p PollManager) stop() {
	p.Terminated <- struct{}{}
}
