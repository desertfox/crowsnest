package job

import (
	"log"
	"time"
)

type SearchService interface {
	Execute(string, string, string, int, []string) ([]byte, error)
	BuildURL(time.Time, time.Time) string
}

type Search struct {
	Type     string        `yaml:"type"`
	Streamid string        `yaml:"streamid"`
	Query    string        `yaml:"query"`
	Fields   []string      `yaml:"fields"`
	Client   SearchService `yaml:"-"`
}

func (s *Search) Run(frequency int) []byte {
	results, err := s.Client.Execute(
		s.Query,
		s.Streamid,
		s.Type,
		frequency,
		s.Fields,
	)
	if err != nil {
		log.Println(err)
		return []byte{}
	}

	return results
}

func (s Search) BuildURL(from, to time.Time) string {
	return s.Client.BuildURL(from, to)
}
