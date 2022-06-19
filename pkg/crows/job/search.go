package job

import (
	"log"
)

type SearchService interface {
	Execute(string, string, int) ([]byte, error)
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
		frequency,
	)
	if err != nil {
		log.Println(err)
		return []byte{}
	}

	return results
}
