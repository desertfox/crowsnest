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

func (j Job) Func(searchService SearchService, reportService ReportService) func() {
	return func() {
		j := j

		log.Println("Job Start, name: " + j.Name)

		count, err := searchService.ExecuteSearch(searchService.GetHeader())
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("Job Results, name: %s, count: %d, alert: %t ", j.Name, count, j.Condition.isAlert(count))

		if j.Output.isVerbose() || j.Condition.isAlert(count) {
			reportService.Send(fmt.Sprintf(`
				🔎 Name  : %s\n\r
				⌚ Freq  : %d\n\r
				📜 Status: %s\n\r
				🧮 Count : %d\n\r
				🔗 Link  : [GrayLog](%s)`,
				j.Name,
				j.Search.Frequency,
				j.Condition.isAlertText(count),
				count,
				searchService.BuildSearchURL(),
			))
		}

		log.Println("Job Finish, name: " + j.Name)
	}
}
