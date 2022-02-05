package job

import (
	"bytes"
	"fmt"
	"log"
	"time"
)

type Condition struct {
	Threshold int    `yaml:"threshold"`
	State     string `yaml:"state"`
}
type Result struct {
	Count int
	When  time.Time
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

func (c *Condition) Parse(rawSearch []byte) Result {
	count := bytes.Count(rawSearch, []byte("\n"))
	if count > 1 {
		count -= 1
	}

	central, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatalf(err.Error())
	}

	return Result{
		Count: count,
		When:  time.Now().In(central),
	}
}

func (r Result) From(f int) time.Time {
	return r.When.Add(time.Duration(-1 * f * int(time.Minute)))
}

func (r Result) To() time.Time {
	return r.When
}
