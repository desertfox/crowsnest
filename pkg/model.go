package crowsnest

import (
	"net/http"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/report"
	"github.com/go-co-op/gocron"
)

var s *gocron.Scheduler = gocron.NewScheduler(time.UTC)

type sessionService interface {
	GetHeader() string
	GetHost() string
}

type queryService interface {
	ExecuteSearch(string) (int, error)
	BuildSearchURL() string
}

type reportService interface {
	Send(string, string, string) error
}

type SearchService struct {
	sessionService
	queryService
}

type Crowsnest struct {
	jobs       []job
	httpClient *http.Client
	sessionService
}

func New(host, username, password, configPath string, httpClient *http.Client) Crowsnest {
	jobs := BuildJobsFromConfig(configPath)

	sessionService, err := session.New(host, username, password, httpClient)
	if err != nil {
		panic(err.Error())
	}

	return Crowsnest{jobs, httpClient, sessionService}
}

func (cn *Crowsnest) ScheduleJobs() {
	for _, job := range cn.jobs {
		query := job.NewSearch(cn.GetHost(), cn.httpClient)

		searchService := SearchService{cn, query}

		s.Every(job.Frequency).Minutes().Do(job.GetCron(searchService, report.Report{}))
	}
}

func (cn Crowsnest) StartBlocking() {
	s.StartBlocking()
}
