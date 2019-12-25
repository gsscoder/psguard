package main

import (
	"flag"
	"fmt"
	"time"
)

const (
	defCPUPercent    = 1
	defMemoryPercent = 1
	defPolling       = time.Second * 2
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
