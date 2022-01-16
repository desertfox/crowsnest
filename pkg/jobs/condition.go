package jobs

import "fmt"

type Condition struct {
	Threshold int    `yaml:"threshold"`
	State     string `yaml:"state"`
}

func (c Condition) isAlert(count int) bool {
	switch c.State {
	case ">":
		return count >= c.Threshold
	case "<":
		return count <= c.Threshold
	default:
		return false
	}
}

func (c Condition) isAlertText(count int) string {
	if c.isAlert(count) {
		return fmt.Sprintf("🔥%d %s= %d🔥", count, c.State, c.Threshold)
	}
	return fmt.Sprintf("✔️%d %s= %d✔️", count, c.State, c.Threshold)
}
