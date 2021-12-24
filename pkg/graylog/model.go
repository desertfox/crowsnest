package graylog

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"
)

type outputService interface {
	Send(string, string) error
}

type Client struct {
	Jobs       []job `yaml:"jobs"`
	lr         *loginRequest
	auth       auth
	httpClient *http.Client
}

type job struct {
	Name      string `yaml:"name"`
	Frequency int    `yaml:"frequency"`
	Option    option `yaml:"options"`
	TeamsURL  string `yaml:"teamsurl"`
	Threshold int    `yaml:"threshold"`
	Type      string `yaml:"type"`
}

type option struct {
	Streamid string   `yaml:"streamid"`
	Query    string   `yaml:"query"`
	Fields   []string `yaml:"fields"`
}

func NewClient(h, u, p string) *Client {
	for i, s := range []string{h, u, p} {
		if s == "" {
			switch i {
			case 0:
				panic("Missing host variable")
			case 1:
				panic("Missing username variable")
			case 2:
				panic("Missing password variable")
			}
		}
	}

	return &Client{[]job{}, &loginRequest{h, u, p}, auth{}, &http.Client{}}
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

func (j job) GetFunc(jobService *Client, outputService outputService) func() {
	return func() {
		fmt.Println("ExecuteJob " + j.Name)

		q := j.newQuery(jobService.lr.Host)

		count, err := q.Execute(jobService.getSessionHeader())
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(time.Now(), count, q.BuildHumanURL())

		var status string
		if count >= j.Threshold {
			status = "ALERT"
		} else {
			status = "OK"
		}

		outputService.Send(j.Name, fmt.Sprintf("Status: %s\nCount: %d\nLink: [GrayLog Query](%s)\n", status, count, q.BuildHumanURL()))
	}
}
