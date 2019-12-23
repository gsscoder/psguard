package main

import (
	"fmt"
	"os/exec"

	"github.com/shirou/gopsutil/process"
)

// ProcessInfo models process relevant informations
type ProcessInfo struct {
	Pid           int32
	Exe           string
	Cmdline       []string
	CPUPercent    float64
	MemoryPercent float64
}

// NewProcessInfoFromExe builds a new ProcessInfo instance for a given name
func NewProcessInfoFromExe(path string) (*ProcessInfo, error) {
	process, err := processFromExe(path)
	if err != nil {
		return nil, err
	}
	return newProcessInfo(*process)
}

// NewProcessInfo builds a new ProcessInfo instance for a given pid
func NewProcessInfo(pid int32) (*ProcessInfo, error) {
	process, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("Process %d not found", pid)
	}
	return newProcessInfo(*process)
}

// Start spawns a process without owning it with original command line
func (p ProcessInfo) Start() error {
	cmd := exec.Command(p.Exe, p.Cmdline...)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Can not start process %d", p.Pid)
	}
	return nil
}

func newProcessInfo(proc process.Process) (*ProcessInfo, error) {
	makeError := func(property string) error {
		return fmt.Errorf("Can not gather %s for process %d", property, proc.Pid)
	}
	exe, err := proc.Exe()
	if err != nil {
		return nil, makeError("executable path")
	}
	cmdline, err := proc.CmdlineSlice()
	if err != nil {
		return nil, makeError("command line")
	}
	cpu, err := proc.CPUPercent()
	if err != nil {
		return nil, makeError("CPU usage")
	}
	mem, err := proc.MemoryPercent()
	if err != nil {
		return nil, makeError("memory usage")
	}
	return &ProcessInfo{
		Pid:           proc.Pid,
		Exe:           exe,
		Cmdline:       cmdline[1:],
		CPUPercent:    cpu,
		MemoryPercent: float64(mem)}, nil
}

func processFromExe(path string) (*process.Process, error) {
	makeError := func() error { return fmt.Errorf("Can not find process with path \"%s\"", path) }
	processes, err := process.Processes()
	if err != nil {
		return nil, makeError()
	}
	for _, process := range processes {
		if exe, _ := process.Exe(); exe == path {
			return process, nil
		}
	}
	return nil, makeError()
}
