package graylog

import (
	"net/http"
)

type Client struct {
	session    *session
	query      Query
	host       string
	httpClient *http.Client
}

func New(un, pw, h string, httpClient *http.Client) *Client {
	return &Client{
		session: newSession(
			h,
			un,
			pw,
			httpClient,
		),
		host:       h,
		httpClient: httpClient,
	}
}

func (c *Client) Execute(query, streamid, typ string, frequency int, fields []string) ([]byte, error) {
	c.query = Query{
		Host:      c.host,
		Query:     query,
		Streamid:  streamid,
		Frequency: frequency,
		Fields:    fields,
		Type:      typ,
	}

	raw, err := c.query.execute(c.session.authHeader(), c.httpClient)
	if err != nil {
		return []byte{}, err
	}
	return raw, nil
}

func (c Client) BuildURL() string {
	return c.query.toURL()
}
