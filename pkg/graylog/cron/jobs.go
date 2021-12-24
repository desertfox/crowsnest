package cron

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"gopkg.in/yaml.v2"
)

type graylogService interface {
	GetSessionHeader() string
	GetHost() string
}

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

func BuildFromConfig(configPath string) *[]job {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err.Error())
	}

	var jobs []job
	err = yaml.Unmarshal(file, &jobs)
	if err != nil {
		panic(err.Error())
	}

	return &jobs
}

func (j job) GetFunc(graylogService graylogService, outputService outputService) func() {
	return func() {
		fmt.Println("ExecuteJob " + j.Name)

		q := j.newQuery(graylogService.GetHost())

		count, err := q.Execute(graylogService.GetSessionHeader())
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
