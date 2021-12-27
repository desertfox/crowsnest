package crowsnest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"gopkg.in/yaml.v2"
)

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

func BuildJobsFromConfig(configPath string) []job {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err.Error())
	}

	var jobs []job
	err = yaml.Unmarshal(file, &jobs)
	if err != nil {
		panic(err.Error())
	}

	return jobs
}

func (j job) GetCron(searchService SearchService, reportService reportService) func() {
	return func() {
		j := j //MARK

		count, err := searchService.ExecuteSearch(searchService.GetHeader())
		if err != nil {
			panic(err.Error())
		}

		reportService.Send(
			j.TeamsURL,
			j.Name,
			fmt.Sprintf("Alert: %s\nCount: %d\nLink: [GrayLog Query](%s)\n", j.shouldAlertText(count), count, searchService.BuildSearchURL()),
		)
	}
}

func (j job) shouldAlertText(count int) string {
	if count >= j.Threshold {
		return fmt.Sprintf("ALERT %d/%d", count, j.Threshold)
	}

	return fmt.Sprintf("OK %d/%d", count, j.Threshold)
}

func (j job) NewSearch(host string, httpClient *http.Client) queryService {
	return search.New(host, j.Name, j.Option.Query, j.Option.Streamid, j.Frequency, j.Option.Fields, httpClient)
}
