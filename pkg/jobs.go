package crowsnest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"gopkg.in/yaml.v2"
)

type job struct {
	Name      string        `yaml:"name"`
	Frequency int           `yaml:"frequency"`
	Threshold int           `yaml:"threshold"`
	TeamsURL  string        `yaml:"teamsurl"`
	Search    searchOptions `yaml:"options"`
}

type searchOptions struct {
	Username string   `yaml:"envusername"`
	Password string   `yaml:"envpassword"`
	Host     string   `yaml:"host"`
	Type     string   `yaml:"type"`
	Streamid string   `yaml:"streamid"`
	Query    string   `yaml:"query"`
	Fields   []string `yaml:"fields"`
}

func (s searchOptions) getUsername() string {
	return os.Getenv(s.Username)
}
func (s searchOptions) getPassword() string {
	return os.Getenv(s.Password)
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

func (j job) NewSearch(httpClient *http.Client) queryService {
	return search.New(j.Name, j.Search.Host, j.Search.Query, j.Search.Streamid, j.Frequency, j.Search.Fields, httpClient)
}

func (j job) GetCron(searchService searchService, reportService reportService) func() {
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
