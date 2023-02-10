package job

import (
	"time"

	"github.com/desertfox/gograylog"
)

type Search struct {
	Type     string   `yaml:"type"`
	Streamid string   `yaml:"streamid"`
	Query    string   `yaml:"query"`
	Fields   []string `yaml:"fields"`
	query    *gograylog.Query
}

func (s *Search) buildQuery(frequency int) {
	s.query = &gograylog.Query{
		QueryString: s.Query,
		StreamID:    s.Streamid,
		Fields:      s.Fields,
		Frequency:   frequency,
		Limit:       10000,
	}
}

func (s Search) BuildURL(host string, from, to time.Time) string {
	return s.query.Url(host, from, to)
}
