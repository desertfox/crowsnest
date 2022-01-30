package job

import (
	"time"
)

const maxHistory int = 10

type History struct {
	results []Result
}

type Result struct {
	Count int
	When  time.Time
}

func newHistory() History {
	return History{
		results: make([]Result, 0, maxHistory),
	}
}

func (h History) Results() []Result {
	return h.results
}

func (h *History) Push(r Result) {
	results := make([]Result, 0, maxHistory)

	results = append(results, r)

	if len(h.results) > 0 {
		results = append(results, h.results[1:]...)
	}

	h.results = results
}
