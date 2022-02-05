package job

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

	h.results = results
}

func (h History) Avg() int {
	if len(h.results) == 0 {
		return 0
	}

	var sum int = 0
	for _, v := range h.results {
		sum += v.Count
	}

	return sum / len(h.results)
}
