package crowsnest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/go-co-op/gocron"
)

var s *gocron.Scheduler = gocron.NewScheduler(time.UTC)

type sessionService interface {
	GetHeader() string
}

type queryService interface {
	ExecuteSearch(string) (int, error)
	BuildSearchURL() string
}

type searchService struct {
	sessionService
	queryService
}

type reportService interface {
	Send(string, string, string) error
}
type crowsnest struct {
	jobs       []Job
	httpClient *http.Client
}

func New(configPath string) crowsnest {
	jobs := BuildJobsFromConfig(configPath)

	for i, job := range jobs {
		color.Yellow(fmt.Sprintf("Loaded Job %d: %s", i, job.Name))
	}

	return crowsnest{jobs, &http.Client{}}
}

func (cn *crowsnest) ScheduleJobs(un, pw string) {
	for i, j := range cn.jobs {
		sessionService := j.NewSession(un, pw, cn.httpClient)

		queryService := j.NewSearch(cn.httpClient)

		searchService := searchService{sessionService, queryService}

		reportService := j.NewReport()

		s.Every(j.Frequency).Minutes().Do(j.GetCron(searchService, reportService))

		color.Green(fmt.Sprintf("Scheduled Job %d: %s", i, j.Name))
	}
}

func (cn crowsnest) StartBlocking() {
	s.StartBlocking()
}
