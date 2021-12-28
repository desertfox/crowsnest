package crowsnest

import (
	"net/http"
	"time"

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
	jobs       []job
	httpClient *http.Client
}

func New(configPath string, httpClient *http.Client) crowsnest {
	jobs := BuildJobsFromConfig(configPath)

	return crowsnest{jobs, httpClient}
}

func (cn *crowsnest) ScheduleJobs() {
	for _, j := range cn.jobs {
		sessionService := j.NewSession(cn.httpClient)

		queryService := j.NewSearch(cn.httpClient)

		searchService := searchService{sessionService, queryService}

		reportService := j.NewReport()

		s.Every(j.Frequency).Minutes().Do(j.GetCron(searchService, reportService))
	}
}

func (cn crowsnest) StartBlocking() {
	s.StartBlocking()
}
