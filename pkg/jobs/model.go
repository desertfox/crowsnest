package jobs

import (
	"errors"
	"fmt"
	"log"
	"os"
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
	Send(string, string) error
}

type Job struct {
	Name      string    `yaml:"name"`
	Condition Condition `yaml:"condition"`
	Output    Output    `yaml:"output"`
	Search    Search    `yaml:"search"`
}

type Condition struct {
	Threshold int    `yaml:"threshold"`
	State     string `yaml:"state"`
}

type Output struct {
	Verbose  int    `yaml:"verbose"`
	TeamsURL string `yaml:"teamsurl"`
}

type Search struct {
	Host      string   `yaml:"host"`
	Type      string   `yaml:"type"`
	Streamid  string   `yaml:"streamid"`
	Query     string   `yaml:"query"`
	Fields    []string `yaml:"fields"`
	Frequency int      `yaml:"frequency"`
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

		output := fmt.Sprintf("Alert: %s\n\rCount: %d\n\rLink: [GrayLog Query](%s)\n\r", j.shouldAlertText(count), count, searchService.BuildSearchURL())

		if j.Output.Verbose > 0 || j.shouldAlert(count) {
			reportService.Send(
				j.Name,
				output,
			)
		}

		if os.Getenv("CROWSNEST_DEBUG") == "log" {
			log.Println(output)
		}

		log.Println("Finished Job: " + j.Name)
	}
}

func (j Job) shouldAlert(count int) bool {
	switch j.Condition.State {
	case ">":
		return count >= j.Condition.Threshold
	case "<":
		return count <= j.Condition.Threshold
	default:
		return false
	}
}

func (j Job) shouldAlertText(count int) string {
	if j.shouldAlert(count) {
		return fmt.Sprintf("ALERT %d/%d", count, j.Condition.Threshold)
	}
	return fmt.Sprintf("OK %d/%d", count, j.Condition.Threshold)
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
