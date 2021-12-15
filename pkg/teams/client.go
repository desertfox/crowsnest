package teams

type Client struct {
	url string
}

func BuildClient(teamsUrl string) *Client {
	return &Client{teamsUrl}
}

func (c Client) Send(title, text string) error {
	return NewReport(title, text).sendReport(c.url)
}
