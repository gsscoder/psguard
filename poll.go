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
	for {
		p.poll()
		time.Sleep(p.Polling)
	}
}

func (p PollManager) poll() {
	for i := 0; i < len(p.Constraints); i++ {
		p.pollGroup(&p.Constraints[i])
	}
}

func (p PollManager) pollGroup(constr *Constraint) {
	for pi := 0; pi < len(constr.Processes); pi++ {
		if constr.Processes[pi].Restarted != nil {
			// Recently restarted, don't need to poll immediately
			switch time.Now().Sub(*constr.Processes[pi].Restarted) < p.Waiting {
			default:
				continue
			case false:
				constr.Processes[pi].Restarted = nil
			}
		}
		switch constr.Processes[pi].Exists() {
		default:
			// Process is running, check resource consumption
			messages := constr.Processes[pi].ViolationMessages(*constr)
			if len(messages) > 0 {
				for _, message := range messages {
					log.Print(message)
				}
			}
		case false:
			log.Printf("Process %s terminated", constr.Name)
			if p.RestartOnEnd {
				// Pocess is terminated and restart is requested
				if err := constr.Processes[pi].Start(); err != nil {
					log.Print(err.Error())
				}
				now := time.Now()
				constr.Processes[pi].Restarted = &now
				log.Printf("Process %s restarted\n", constr.Processes[pi].Exe)
			}
		}
	}
}
