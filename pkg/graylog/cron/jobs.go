package cron

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type sessionService interface {
	GetHeader() string
	GetHost() string
}

type outputService interface {
	Send(string, string) error
}

type Riddler interface {
	Execute(string) (int, error)
	BuildHumanURL() string
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

func (j job) GetFunc(sessionService sessionService, outputService outputService, q Riddler) func() {
	return func() {
		j := j //MARK
		count, err := q.Execute(sessionService.GetHeader())
		if err != nil {
			panic(err.Error())
		}

		outputService.Send(j.Name, fmt.Sprintf("Alert: %s\nCount: %d\nLink: [GrayLog Query](%s)\n", j.shouldAlertText(count), count, q.BuildHumanURL()))
	}
}

func (j job) shouldAlertText(count int) string {
	if count >= j.Threshold {
		return fmt.Sprintf("ALERT %d/%d", count, j.Threshold)
	}

	return fmt.Sprintf("OK %d/%d", count, j.Threshold)
}
