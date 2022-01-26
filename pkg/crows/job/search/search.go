package search

import (
	"log"
)

type SearchService interface {
	Execute(string, string, string, int, []string) ([]byte, error)
	BuildURL() string
}

type Search struct {
	Type      string        `yaml:"type"`
	Streamid  string        `yaml:"streamid"`
	Query     string        `yaml:"query"`
	Fields    []string      `yaml:"fields"`
	Condition Condition     `yaml:"condition"`
	Output    Output        `yaml:"output"`
	Client    SearchService `yaml:"-"`
}

func (s *Search) Run(frequency int) {
	results, err := s.Client.Execute(
		s.Query,
		s.Streamid,
		s.Type,
		frequency,
		s.Fields,
	)
	if err != nil {
		log.Fatal(err)
	}

	s.Condition.Measure(results)
}

func (s Search) BuildURL() string {
	return s.Client.BuildURL()
}
