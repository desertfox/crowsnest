package job

import (
	"fmt"
	"log"

	"github.com/desertfox/crowsnest/pkg/condition"
	"github.com/desertfox/crowsnest/pkg/output"
	"github.com/desertfox/crowsnest/pkg/search"
)

type Job struct {
	Name      string              `yaml:"name"`
	Condition condition.Condition `yaml:"condition"`
	Output    output.Output       `yaml:"output"`
	Search    search.Search       `yaml:"search"`
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

func (j Job) Func(searchService search.SearchService, reportService output.ReportService) func() {
	return func() {
		j := j

		log.Println("Job Start, name: " + j.Name)

		count, err := searchService.ExecuteSearch(searchService.GetHeader())
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("Job Results, name: %s, count: %d, alert: %t ", j.Name, count, j.Condition.IsAlert(count))

		if j.Output.IsVerbose() || j.Condition.IsAlert(count) {
			reportService.Send(
				fmt.Sprintf("🔎 Name  : %s\n\r"+
					"⌚ Freq  : %d\n\r"+
					"📜 Status: %s\n\r"+
					"🧮 Count : %d\n\r"+
					"🔗 Link  : [GrayLog](%s)",
					j.Name,
					j.Search.Frequency,
					j.Condition.IsAlertText(count),
					count,
					searchService.BuildSearchURL(),
				),
			)
		}

		log.Println("Job Finish, name: " + j.Name)
	}
}
