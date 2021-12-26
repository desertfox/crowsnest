package cron

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type graylogService interface {
	BuildSearchURL() string
	ExecuteSearch() (int, error)
}

type reportService interface {
	Send(string, string, string) error
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

func (j job) GetCron(graylogService graylogService, reportService reportService) func() {
	return func() {
		j := j //MARK

		count, err := graylogService.ExecuteSearch()
		if err != nil {
			panic(err.Error())
		}

		reportService.Send(
			j.TeamsURL,
			j.Name,
			fmt.Sprintf("Alert: %s\nCount: %d\nLink: [GrayLog Query](%s)\n", j.shouldAlertText(count), count, graylogService.BuildSearchURL()),
		)
	}
}

func (j job) shouldAlertText(count int) string {
	if count >= j.Threshold {
		return fmt.Sprintf("ALERT %d/%d", count, j.Threshold)
	}

	return fmt.Sprintf("OK %d/%d", count, j.Threshold)
}
