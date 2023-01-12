package job

import (
	"fmt"
)

type Condition struct {
	Threshold int    `yaml:"threshold"`
	State     string `yaml:"state"`
}

func (c Condition) IsAlert(r Result) bool {
	switch c.State {
	case ">":
		return r.Count >= c.Threshold
	case "<":
		return r.Count <= c.Threshold
	default:
		return false
	}
}

func (c Condition) IsAlertText(r Result) string {
	if c.IsAlert(r) {
		return fmt.Sprintf("🔥%d %s= %d🔥", r.Count, c.State, c.Threshold)
	}
	return fmt.Sprintf("✔️%d %s= %d✔️", r.Count, c.State, c.Threshold)
}
