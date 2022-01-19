package graylog

import (
	"net/http"
	"strings"
)

type Graylog struct {
	session    *session
	query      query
	httpClient *http.Client
}

func New(un, pw, h, q, streamid string, frequency int, fields []string, t string, httpClient *http.Client) *Graylog {
	return &Graylog{
		session: newSession(
			h,
			un,
			pw,
			httpClient,
		),
		query: query{
			h,
			q,
			streamid,
			frequency,
			fields,
			t,
		},
		httpClient: httpClient,
	}
}

func (g Graylog) Execute() (int, error) {
	raw, err := g.query.execute(g.session.authHeader(), g.httpClient)
	if err != nil {
		return 0, err
	}

	count := strings.Count(string(raw), "\n")
	if count > 1 {
		count -= 1
	}

	return count, nil
}

func (g Graylog) BuildURL() string {
	return g.query.toURL()
}
