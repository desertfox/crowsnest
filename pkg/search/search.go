package search

import (
	"net/http"

	"github.com/desertfox/crowsnest/pkg/graylog/query"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
)

type SessionService interface {
	GetHeader() string
}

type QueryService interface {
	ExecuteSearch(string) (int, error)
	BuildSearchURL() string
}

type SearchService struct {
	SessionService
	QueryService
}

type Search struct {
	Host      string   `yaml:"host"`
	Type      string   `yaml:"type"`
	Streamid  string   `yaml:"streamid"`
	Query     string   `yaml:"query"`
	Fields    []string `yaml:"fields"`
	Frequency int      `yaml:"frequency"`
}

func (s Search) SearchService(un, pw string, httpClient *http.Client) SearchService {
	return SearchService{
		SessionService: session.New(
			s.Host,
			un,
			pw,
			httpClient,
		),
		QueryService: query.New(
			s.Host,
			s.Query,
			s.Streamid,
			s.Frequency,
			s.Fields,
			s.Type,
			httpClient,
		),
	}
}
