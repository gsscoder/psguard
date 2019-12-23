package main

import (
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
			constraints[i].CPUPercent = defCPUPercent
			log.Printf("CPU constraint for %s invalid, setted to default (%.2f%%)", constr.Name, options.CPUPercent)
		}
		if constr.MemoryPercent < 0 {
			constraints[i].MemoryPercent = defMemoryPercent
			log.Printf("Memory constraint for %s invalid, setted to default (%.2f%%)", constr.Name, options.MemoryPercent)
		}
	}
	server := NewPollManager(*options, constraints)
	server.Start()
}
