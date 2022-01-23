package graylog

import (
	"net/http"
)

type Client struct {
	session    *session
	query      query
	httpClient *http.Client
}

func New(un, pw, h, q, streamid string, frequency int, fields []string, t string, httpClient *http.Client) *Client {
	return &Client{
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

func (c Client) Execute() ([]byte, error) {
	raw, err := c.query.execute(c.session.authHeader(), c.httpClient)
	if err != nil {
		return []byte{}, err
	}
	return raw, nil
}

func (c Client) BuildURL() string {
	return c.query.toURL()
}
