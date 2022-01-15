package jobs

import (
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

func (j Job) GetCron(searchService SearchService, reportService ReportService) func() {
	return func() {
		j := j //MARK

		log.Println("Job Start: " + j.Name)

		count, err := searchService.ExecuteSearch(searchService.GetHeader())
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("Job %s Results count: %d, alert: %t ", j.Name, count, j.shouldAlert(count))

		output := fmt.Sprintf("⌚ Freq  : %d\n\r", j.Search.Frequency)
		output += fmt.Sprintf("📜 Status: %s\n\r", j.shouldAlertText(count))
		output += fmt.Sprintf("🧮 Count : %d\n\r", count)
		output += fmt.Sprintf("🔗 Link  : [GrayLog Query](%s)\n\r", searchService.BuildSearchURL())

		if j.Output.Verbose > 0 || j.shouldAlert(count) {
			reportService.Send(
				"🔎 Name: "+j.Name,
				output,
			)
		}

		log.Println("Job Finish: " + j.Name)
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
		return fmt.Sprintf("🔥%d %s= %d🔥", count, j.Condition.State, j.Condition.Threshold)
	}
	return fmt.Sprintf("✔️%d %s= %d✔️", count, j.Condition.State, j.Condition.Threshold)
}
