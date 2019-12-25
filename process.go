package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"github.com/shirou/gopsutil/process"
)

// ProcessInfo models process relevant informations
type ProcessInfo struct {
	Pid           int32
	Exe           string
	Cmdline       []string
	CPUPercent    float64
	MemoryPercent float64
	Restarted     *time.Time
}

// ProcessInfoSliceFromRegexp builds a new slice of ProcessInfo instances for a pattern
func ProcessInfoSliceFromRegexp(re regexp.Regexp) []ProcessInfo {
	result := make([]ProcessInfo, 0)
	processes, err := process.Processes()
	if err != nil {
		return result
	}
	for _, process := range processes {
		exe, err := process.Exe()
		if err == nil {
			if re.MatchString(exe) {
				info, err := newProcessInfo(*process)
				if err == nil {
					result = append(result, *info)
				}
			}
		}
	}
	return result
}

// Exists checks if a process does exists
func (p ProcessInfo) Exists() bool {
	exists, _ := process.PidExists(p.Pid)
	return exists
}

func (p ProcessInfo) ViolationMessages(constr Constraint) []string {
	messages := make([]string, 0)
	if p.CPUPercent > constr.CPUPercent {
		messages = append(messages,
			fmt.Sprintf("%s: CPU constraint of %.2f%% violated by +%.2f%%",
				constr.Name, constr.CPUPercent, p.CPUPercent-constr.CPUPercent))
	}
	if p.MemoryPercent > constr.MemoryPercent {
		messages = append(messages,
			fmt.Sprintf("%s: Memory constraint of %.2f%% violated by +%.2f%%",
				constr.Name, constr.MemoryPercent, p.MemoryPercent-constr.MemoryPercent))
	}
	return messages
}

// Start spawns a process without owning it with original command line
func (p *ProcessInfo) Start() error {
	cmd := exec.Command(p.Exe, p.Cmdline...)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Can not start process %d", p.Pid)
	}
	p.Pid = int32(cmd.Process.Pid)
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
