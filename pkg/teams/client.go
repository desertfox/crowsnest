package teams

type Client struct {
	url string
}

func BuildClient(teamsUrl string) *Client {
	return &Client{teamsUrl}
}

func (c Client) Send(text string) error {
	return NewReport("", text).sendReport(c.url)
}
