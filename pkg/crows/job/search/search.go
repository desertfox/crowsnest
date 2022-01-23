package search

import (
	"log"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/crows/job/search/graylog"
)

type SearchService interface {
	Execute() ([]byte, error)
	BuildURL() string
}

type Search struct {
	Type      string    `yaml:"type"`
	Streamid  string    `yaml:"streamid"`
	Query     string    `yaml:"query"`
	Fields    []string  `yaml:"fields"`
	Condition Condition `yaml:"condition"`
	client    SearchService
}

func (s *Search) Run(host, un, pw string, frequency int, httpClient *http.Client) {
	client := graylog.New(
		un,
		pw,
		host,
		s.Query,
		s.Streamid,
		frequency,
		s.Fields,
		s.Type,
		httpClient,
	)

	s.client = client

	results, err := s.client.Execute()
	if err != nil {
		log.Fatal(err)
	}

	s.Condition.Eval(results)
}

func (s Search) BuildURL() string {
	return s.client.BuildURL()
}
