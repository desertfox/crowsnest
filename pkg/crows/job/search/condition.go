package search

import (
	"bytes"
	"fmt"
)

type Condition struct {
	Threshold int    `yaml:"threshold"`
	State     string `yaml:"state"`
	Count     int
}

func (c Condition) IsAlert() bool {
	switch c.State {
	case ">":
		return c.Count >= c.Threshold
	case "<":
		return c.Count <= c.Threshold
	default:
		return false
	}
}

func (c Condition) IsAlertText() string {
	if c.IsAlert() {
		return fmt.Sprintf("🔥%d %s= %d🔥", c.Count, c.State, c.Threshold)
	}
	return fmt.Sprintf("✔️%d %s= %d✔️", c.Count, c.State, c.Threshold)
}

func (c *Condition) Eval(rawSearch []byte) {
	count := bytes.Count(rawSearch, []byte("\n"))
	if count > 1 {
		count -= 1
	}
	c.Count = count
}
