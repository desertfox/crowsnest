package job

import "log"

const maxHistory int = 10

type History struct {
	results []Result
}

func newHistory() *History {
	return &History{
		results: []Result{},
	}
}

func (h History) Results() []Result {
	return h.results
}

func (h *History) Push(r Result) {
	results := []Result{r}

	results = append(results, h.results...)

	if len(h.results) >= maxHistory {
		results = results[:maxHistory]
	}

	log.Printf("results %#v", results)

	h.results = results
}
