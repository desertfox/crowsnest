package condition

import "fmt"

type Condition struct {
	Threshold int    `yaml:"threshold"`
	State     string `yaml:"state"`
}

func (c Condition) IsAlert(count int) bool {
	switch c.State {
	case ">":
		return count >= c.Threshold
	case "<":
		return count <= c.Threshold
	default:
		return false
	}
}

func (c Condition) IsAlertText(count int) string {
	if c.IsAlert(count) {
		return fmt.Sprintf("🔥%d %s= %d🔥", count, c.State, c.Threshold)
	}
	return fmt.Sprintf("✔️%d %s= %d✔️", count, c.State, c.Threshold)
}
