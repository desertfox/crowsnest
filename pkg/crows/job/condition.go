package job

import (
	"fmt"
)

type Condition struct {
	Threshold int    `yaml:"threshold"`
	Operator  string `yaml:"operator"`
}

// Takes a Result and operates using configured threshold and state condition attributes.
func (c Condition) IsAlert(count int) bool {
	switch c.Operator {
	case ">":
		return count >= c.Threshold
	case "<":
		return count <= c.Threshold
	default:
		return false
	}
}

func (c Condition) IsAlertText(alert bool, count int) string {
	if alert {
		return fmt.Sprintf("🔥%d %s= %d🔥", count, c.Operator, c.Threshold)
	}
	return fmt.Sprintf("✔️%d %s= %d✔️", count, c.Operator, c.Threshold)
}
