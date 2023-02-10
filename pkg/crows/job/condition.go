package job

import (
	"fmt"
)

type Condition struct {
	Threshold int    `yaml:"threshold"`
	State     string `yaml:"state"`
}

// Takes a Result and operates using configured threshold and state condition attributes.
func (c Condition) IsAlert(r *Result) {
	switch c.State {
	case ">":
		r.Alert = r.Count >= c.Threshold
	case "<":
		r.Alert = r.Count <= c.Threshold
	default:
		r.Alert = false
	}
}

func (c Condition) IsAlertText(r *Result) string {
	if r.Alert {
		return fmt.Sprintf("🔥%d %s= %d🔥", r.Count, c.State, c.Threshold)
	}
	return fmt.Sprintf("✔️%d %s= %d✔️", r.Count, c.State, c.Threshold)
}
