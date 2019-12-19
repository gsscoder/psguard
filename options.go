package main

import (
	"flag"
	"fmt"
	"time"
)

const (
	defCPUPercent    = 1
	defMemoryPercent = 1
	defPolling       = time.Second
	defWaiting       = time.Second * 5
)

// Options models command line options
type Options struct {
	Pid           int64
	CPUPercent    float64
	MemoryPercent float64
	Polling       time.Duration
	RestartOnEnd  bool
	Waiting       time.Duration
}

// NewOptions parses command line options
func NewOptions() *Options {
	opts := Options{}
	var version bool
	flag.BoolVar(&version, "version", false, "displays version information")
	flag.Int64Var(&opts.Pid, "pid", 0, "pid of the process to monitor")
	flag.Float64Var(&opts.CPUPercent, "cpu", defCPUPercent, "allowed CPU % usage")
	flag.Float64Var(&opts.MemoryPercent, "mem", defMemoryPercent, "allowed memory % usage")
	flag.DurationVar(&opts.Polling, "poll", defPolling, "defines polling interval")
	flag.BoolVar(&opts.RestartOnEnd, "restart", false, "restart process if terminated")
	flag.DurationVar(&opts.Waiting, "wait", defWaiting, "time to wait before polling again after a restart")
	flag.Usage = func() {
		PrintInfo()
		fmt.Println("usage:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if version {
		return nil
	}
	return &opts
}

// Validate validates command line options
func (o Options) Validate() error {
	if o.Pid <= 0 {
		return fmt.Errorf("-pid is mandatory")
	}
	if o.CPUPercent < 0 || o.CPUPercent > 100 {
		return fmt.Errorf("-cpu must be between 0 and 100")
	}
	if o.MemoryPercent < 0 || o.MemoryPercent > 100 {
		return fmt.Errorf("-mem must be between 0 and 100")
	}
	return nil
}
