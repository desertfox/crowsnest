package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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
	if os.Getenv("CROWS_DEBUG") != "" {
		fmt.Printf("DEBUG: rawSearch %s\n", rawSearch)

	}

	count := bytes.Count(rawSearch, []byte("\n"))

	//BUG: API sometimes returns complex object?
	if count == 0 && len(rawSearch) > 0 {
		var hack map[string]interface{}
		err := json.Unmarshal(rawSearch, &hack)
		if err != nil {
			fmt.Printf("Error parsing json. data:%s\nerr:%s\n", rawSearch, err)
		}

		if val, ok := hack["total_results"]; ok {
			return Result{
				Count: int(val.(float64)),
				When:  time.Now(),
			}
		}
	}

	//Remove csv headers
	if count > 2 {
		count -= 1
	}

	return Result{
		Count: count,
		When:  time.Now(),
	}
}

func (r Result) From(f int) time.Time {
	return r.When.Add(time.Duration(-1 * f * int(time.Minute)))
}

func (r Result) To() time.Time {
	return r.When
}
