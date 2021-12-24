package graylog

import (
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

type sessionService interface {
	GetSessionHeader(*http.Client) string
	GetHost() string
}

type Client struct {
	Jobs       []job `yaml:"jobs"`
	lr         sessionService
	httpClient *http.Client
}

func NewClient(s sessionService) *Client {
	return &Client{[]job{}, s, &http.Client{}}
}

func (c *Client) BuildJobsFromConfig(configPath string) *Client {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err.Error())
	}

	err = yaml.Unmarshal(file, &c.Jobs)
	if err != nil {
		panic(err.Error())
	}

	return c
}

func (c *Client) getSessionHeader() string {
	return c.lr.GetSessionHeader(c.httpClient)
}
