package job

import (
	"net/http"

	"github.com/desertfox/crowsnest/pkg/graylog"
)

type SearchService interface {
	Execute() ([]byte, error)
	BuildURL() string
}

type Search struct {
	Host      string   `yaml:"host"`
	Type      string   `yaml:"type"`
	Streamid  string   `yaml:"streamid"`
	Query     string   `yaml:"query"`
	Fields    []string `yaml:"fields"`
	Frequency int      `yaml:"frequency"`
}

func (s Search) Service(un, pw string, httpClient *http.Client) SearchService {
	return graylog.New(
		un,
		pw,
		s.Host,
		s.Query,
		s.Streamid,
		s.Frequency,
		s.Fields,
		s.Type,
		httpClient,
	)
}
