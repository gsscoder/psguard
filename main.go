package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	options := NewOptions()
	if options == nil {
		PrintInfo()
		os.Exit(0)
	}
	constraints := LoadConstraintList()
	if len(constraints) == 0 {
		Fail("Process list is empty")
	}
	for i, constr := range constraints {
		if constr.CPUPercent < 0 {
			Fail(fmt.Sprintf("Invalid CPU constraint for %s", constr.Name))
		}
		if constr.MemoryPercent < 0 {
			Fail(fmt.Sprintf("Invalid Memory constraint for %s", constr.Name))
		}
		process, err := NewProcessInfo(constr.Pid)
		if err == nil {
			constraints[i].Process = process
		} else {
			log.Fatalf("Can't find process %s", constr.Name)
		}
	}
	boundConstraints := ExcludeUnboundPids(constraints)
	if len(boundConstraints) == 0 {
		Fail("No process to poll")
	}

	server := NewPollManager(*options, boundConstraints)
	server.Start()
}
