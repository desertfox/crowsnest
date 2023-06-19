package job

import "time"

const maxHistory int = 20

// History container for Results
type History struct {
	Results    []*Result `json:"results"`
	AlertCount int       `json:"alertCount"`
}

type Result struct {
	Count int       `json:"Count"`
	When  time.Time `json:"When"`
	Alert bool      `json:"Alert"`
}

func newHistory() *History {
	return &History{
		Results:    make([]*Result, 0, maxHistory),
		AlertCount: 0,
	}
}

// Add
func (h *History) Add(r *Result) {
	if r.Alert {
		h.AlertCount++
	} else {
		h.AlertCount = 0
	}

	results := []*Result{r}
	results = append(results, h.Results...)

	if len(h.Results) >= maxHistory {
		results = results[:maxHistory]
	}

	h.Results = results
}

func (h History) Avg() int {
	if len(h.Results) == 0 {
		return 0
	}

	sum := 0
	for i := 0; i < len(h.Results); i++ {
		sum += h.Results[i].Count
	}

	return sum / len(h.Results)
}

func (r Result) From(f int) time.Time {
	return r.When.Add(time.Duration(int64(-1) * int64(f) * int64(time.Minute)))
}

func (r Result) To() time.Time {
	return r.When
}
