package graylog

import "net/http"

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
	return g.query.execute(g.session.authHeader(), g.httpClient)
}

func (g Graylog) BuildURL() string {
	return g.query.toURL()
}
