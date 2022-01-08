package jobs

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

type SessionService interface {
	GetHeader() string
}

type QueryService interface {
	ExecuteSearch(string) (int, error)
	BuildSearchURL() string
}

type SearchService struct {
	SessionService
	QueryService
}

type ReportService interface {
	Send(string, string, string) error
}

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

type JobList []Job

func NewJob(n string, f, t int, teamurl string, so SearchOptions) Job {
	return Job{n, f, t, teamurl, so}
}

func NewSearchOptions(h, t, streamid, q string, fields []string, from, to string) SearchOptions {
	return SearchOptions{h, t, streamid, q, fields, from, to}
}

func (j Job) GetCron(searchService SearchService, reportService ReportService) func() {
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

func (jl JobList) checkIfExists(j Job) bool {
	for _, job := range jl {
		if job.Name == j.Name {
			return true
		}
	}
	return false
}

func (jl *JobList) Add(j Job) error {
	if jl.checkIfExists(j) {
		return errors.New("job exists")
	}

	*jl = append(*jl, j)

	return nil
}
