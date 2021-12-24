package graylog

import (
	"fmt"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog/search"
)

type outputService interface {
	Send(string, string) error
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

func (j job) GetFunc(jobService *Client, outputService outputService) func() {
	return func() {
		fmt.Println("ExecuteJob " + j.Name)

		q := j.newQuery(jobService.lr.GetHost())

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

func (j job) newQuery(host string) search.Query {
	return search.NewQuery(host, j.Name, j.Option.Query, j.Option.Streamid, j.Frequency, j.Option.Fields)
}
