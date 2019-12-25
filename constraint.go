package main

import (
	"fmt"
	"regexp"

	"github.com/tidwall/gjson"
)

type Constraint struct {
	Name          string
	CPUPercent    float64
	MemoryPercent float64
	Processes     []ProcessInfo
}

type ConstraintList []Constraint

// NewConstraintList builds a list of constraint from configuration file
func NewConstraintList() ConstraintList {
	config := ReadText(ConfigPath)
	if len(config) == 0 {
		return []Constraint{}
	}
	result := gjson.Get(config, "constraints.process\\.groups")
	constraints := make([]Constraint, 0)
	result.ForEach(func(key, value gjson.Result) bool {
		name := key.Str
		item := value.Raw
		matches := gjson.Get(item, "match").Array()
		processes := make([]ProcessInfo, 0)
		for _, match := range matches {
			re, err := regexp.Compile(match.Str)
			if err == nil {
				processes = append(processes,
					ProcessInfoSliceFromRegexp(*re)...)
			}
		}
		cpu := gjson.Get(item, "cpu").Num
		mem := gjson.Get(item, "mem").Num
		constraints = append(constraints, Constraint{
			Name:          name,
			CPUPercent:    cpu,
			MemoryPercent: mem,
			Processes:     processes,
		})
		return true
	})
	return constraints
}

// Sanitize checks and sanitizes constraint values and produces messages when violated
func (c ConstraintList) Sanitize(options Options) []string {
	messages := make([]string, 0)
	for i := 0; i < len(c); i++ {
		if c[i].CPUPercent < 0 {
			c[i].CPUPercent = defCPUPercent
			messages = append(messages,
				fmt.Sprintf("CPU constraint for %s invalid, setted to default (%.2f%%)",
					c[i].Name, options.CPUPercent))
		}
		if c[i].MemoryPercent < 0 {
			c[i].MemoryPercent = defMemoryPercent
			messages = append(messages,
				fmt.Sprintf("Memory constraint for %s invalid, setted to default (%.2f%%)",
					c[i].Name, options.MemoryPercent))
		}
	}
	return messages
}

// HasProcesses checks if at least one process is available for polling
func (c ConstraintList) HasProcesses() bool {
	for _, constr := range c {
		if len(constr.Processes) > 0 {
			return true
		}
	}
	return false
}
