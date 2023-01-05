package job

import (
	"log"
	"time"

	"github.com/desertfox/gograylog"
)

type Search struct {
	Type     string            `yaml:"type"`
	Streamid string            `yaml:"streamid"`
	Query    string            `yaml:"query"`
	Fields   []string          `yaml:"fields"`
	Client   *gograylog.Client `yaml:"-"`
}

func (s *Search) Run(frequency int) []byte {

	q := gograylog.Query{
		QueryString: s.Query,
		StreamID:    s.Streamid,
		Fields:      s.Fields,
		Frequency:   frequency,
		Limit:       10000,
	}

	results, err := s.Client.Search(q)
	if err != nil {
		log.Println(err)
		return []byte{}
	}

	return results
}

func (s Search) BuildURL(from, to time.Time) string {
	q := gograylog.Query{
		QueryString: s.Query,
		StreamID:    s.Streamid,
		Fields:      s.Fields,
		Limit:       10000,
	}

	return q.Url(s.Client.Host, from, to)
}
