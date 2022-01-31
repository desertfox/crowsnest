package job

import "log"

const maxHistory int = 10

type History struct {
	results []Result
}

func newHistory() *History {
	return &History{
		results: make([]Result, 0, maxHistory),
	}
}

func (h History) Results() []Result {
	return h.results
}

func (h *History) Push(r Result) {
	results := make([]Result, 0, maxHistory)

	log.Printf("Empty %#v", results)

	results = append(results, r)

	log.Printf("NewOne %#v", results)

	log.Printf("Current %#v", h.results[1:])

	if len(h.results) > 0 {
		results = append(results, h.results[1:]...)
	}

	log.Printf("The rest %#v", results)

	h.results = results
}
