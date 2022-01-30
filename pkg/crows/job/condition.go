package job

import (
	"bytes"
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

func (c *Condition) Parse(rawSearch []byte) int {
	count := bytes.Count(rawSearch, []byte("\n"))
	if count > 1 {
		count -= 1
	}
	return count
}
