package graylog

type sessionService interface {
	GetHeader() string
	GetHost() string
}

type Client struct {
	sessionService sessionService
}

func NewClient(s sessionService) *Client {
	return &Client{s}
}

func (c *Client) GetSessionHeader() string {
	return c.sessionService.GetHeader()
}

func (c *Client) GetHost() string {
	return c.sessionService.GetHost()
}
