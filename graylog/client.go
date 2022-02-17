package graylog

import (
	"net/http"
	"time"
)

var HttpClient = http.Client{}

type Client struct {
	session *session
	query   Query
}

func New(h string) *Client {
	return &Client{
		session: newSession(
			h,
		),
	}
}

func (c *Client) Execute(query, streamid, typ string, frequency int, fields []string) ([]byte, error) {
	c.query = Query{
		Host:      c.session.loginRequest.Host,
		Query:     query,
		Streamid:  streamid,
		Frequency: frequency,
		Fields:    fields,
		Type:      typ,
	}

	raw, err := c.query.execute(c.session.authHeader())
	if err != nil {
		return []byte{}, err
	}
	return raw, nil
}

func (c Client) BuildURL(from, to time.Time) string {
	return c.query.toURL(from, to)
}
