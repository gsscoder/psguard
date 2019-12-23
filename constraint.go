package main

import (
	"strings"

	"github.com/tidwall/gjson"
)

type Constraint struct {
	Name          string
	Pid           int32
	CPUPercent    float64
	MemoryPercent float64
	Process       *ProcessInfo
}

func LoadConstraintList() []Constraint {
	config := strings.TrimSpace(ReadText(ConfigPath))
	if len(config) == 0 {
		return []Constraint{}
	}
	result := gjson.Get(config, "constraints")
	constraints := make([]Constraint, 0)
	result.ForEach(func(key, value gjson.Result) bool {
		name := key.String()
		item := value.String()
		pid := gjson.Get(item, "pid").Int()
		cpu := gjson.Get(item, "cpu").Float()
		mem := gjson.Get(item, "mem").Float()
		constraints = append(constraints, Constraint{
			Name:          name,
			Pid:           int32(pid),
			CPUPercent:    cpu,
			MemoryPercent: mem,
		})
		return true
	})
	return constraints
}

func ExcludeUnboundPids(constraints []Constraint) []Constraint {
	filtered := make([]Constraint, 0)
	for _, constr := range constraints {
		if constr.Process != nil {
			filtered = append(filtered, constr)
		}
	}
	return filtered
}
