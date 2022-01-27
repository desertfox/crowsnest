package search

import (
	"fmt"
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

func (s *Search) Send(name string, frequency int) {
	if s.Output.IsVerbose() || s.Condition.IsAlert() {
		s.Output.Send(
			s.Output.URL,
			s.buildText(name, frequency),
		)
	}
}

func (s *Search) buildText(name string, frequency int) string {
	return fmt.Sprintf("🔎 Name  : %s\n\r"+
		"⌚ Freq  : %d\n\r"+
		"📜 Status: %s\n\r"+
		"🧮 Count : %d\n\r"+
		"🔗 Link  : [GrayLog](%s)",
		name,
		frequency,
		s.Condition.IsAlertText(),
		s.Condition.Count,
		s.BuildURL(),
	)
}

func (s Search) BuildURL() string {
	return s.Client.BuildURL()
}
