package crowsnest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/teams/report"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type Job struct {
	Name      string        `yaml:"name"`
	Frequency int           `yaml:"frequency"`
	Threshold int           `yaml:"threshold"`
	TeamsURL  string        `yaml:"teamsurl"`
	Search    SearchOptions `yaml:"options"`
}

type SearchOptions struct {
	Host     string   `yaml:"host"`
	Type     string   `yaml:"type"`
	Streamid string   `yaml:"streamid"`
	Query    string   `yaml:"query"`
	Fields   []string `yaml:"fields"`
	From     string   `yaml:"from"`
	To       string   `yaml:"to"`
}

func NewJob(n string, f, t int, teamurl string, so SearchOptions) Job {
	return Job{n, f, t, teamurl, so}
}

func NewSearchOptions(h, t, streamid, q string, fields []string, from, to string) SearchOptions {
	return SearchOptions{h, t, streamid, q, fields, from, to}

}

func BuildJobsFromConfig(configPath string) []Job {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err.Error())
	}

	data := make(map[string][]Job)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		panic(err.Error())
	}

	return data["jobs"]
}

func (j Job) NewSession(un, pw string, httpClient *http.Client) sessionService {
	return session.New(j.Search.Host, un, pw, httpClient)
}

func (j Job) NewSearch(httpClient *http.Client) queryService {
	return search.New(
		j.Search.Host,
		j.Search.Query,
		j.Search.Streamid,
		j.Frequency,
		j.Search.Fields,
		httpClient,
	)
}

func (j Job) NewReport() reportService {
	return report.Report{
		Url: j.TeamsURL,
	}
}

func (j Job) GetCron(searchService searchService, reportService reportService) func() {
	return func() {
		j := j //MARK

		color.Yellow("Starting Job: " + j.Name)

		count, err := searchService.ExecuteSearch(searchService.GetHeader())
		if err != nil {
			panic(err.Error())
		}

		color.Blue("Search Complete: " + j.Name)

		reportService.Send(
			j.Name,
			searchService.BuildSearchURL(),
			fmt.Sprintf("Alert: %s\nCount: %d\nLink: [GrayLog Query](%s)\n", j.shouldAlertText(count), count, searchService.BuildSearchURL()),
		)

		color.Green("Finished Job: " + j.Name)
	}
}

func (j Job) shouldAlertText(count int) string {
	if count >= j.Threshold {
		return fmt.Sprintf("ALERT %d/%d", count, j.Threshold)
	}

	return fmt.Sprintf("OK %d/%d", count, j.Threshold)
}
