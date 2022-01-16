package jobs

import (
	"fmt"
	"log"
)

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
	Send(string) error
}

func (o Output) isVerbose() bool {
	return o.Verbose > 0
}

func (j Job) Func(searchService SearchService, reportService ReportService) func() {
	return func() {
		j := j //MARK

		log.Println("Job Start: " + j.Name)

		count, err := searchService.ExecuteSearch(searchService.GetHeader())
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("Job %s Results count: %d, alert: %t ", j.Name, count, j.shouldAlert(count))

		output := fmt.Sprintf(`
		🔎 Name  : %s\n\r
		⌚ Freq  : %d\n\r
		📜 Status: %s\n\r
		🧮 Count : %d\n\r
		🔗 Link  : [GrayLog](%s)`,
			j.Name,
			j.Search.Frequency,
			j.shouldAlertText(count),
			count,
			searchService.BuildSearchURL(),
		)

		if j.Output.isVerbose() || j.shouldAlert(count) {
			reportService.Send(output)
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
