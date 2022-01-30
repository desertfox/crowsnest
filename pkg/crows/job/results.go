package job

import "time"

type Results []Result

type Result struct {
	Count int
	When  time.Time
}
