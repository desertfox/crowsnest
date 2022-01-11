package jobs

import (
	"errors"
	"fmt"
	"log"
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
	Verbose   int           `yaml:"verbose"`
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

func (j Job) GetCron(searchService SearchService, reportService ReportService) func() {
	return func() {
		j := j //MARK

		log.Println("Starting Job: " + j.Name)

		count, err := searchService.ExecuteSearch(searchService.GetHeader())
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("Search Complete: " + j.Name)

		output := fmt.Sprintf("Alert: %s\nCount: %d\nLink: [GrayLog Query](%s)\n", j.shouldAlertText(count), count, searchService.BuildSearchURL())

		if j.Verbose > 0 || j.shouldAlert(count) {
			reportService.Send(
				j.Name,
				searchService.BuildSearchURL(),
				output,
			)
		} else {
			log.Println(output)
		}

		log.Println("Finished Job: " + j.Name)
	}
}

func (j Job) shouldAlert(count int) bool {
	return count >= j.Threshold
}

func (j Job) shouldAlertText(count int) string {
	if j.shouldAlert(count) {
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

//Add Job to JobList if j.Name does not already exist.
func (jl *JobList) Add(j Job) error {
	if jl.checkIfExists(j) {
		return errors.New("job exists")
	}

	*jl = append(*jl, j)

	return nil
}
