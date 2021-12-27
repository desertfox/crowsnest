package crowsnest

import (
	"net/http"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/teams/report"
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

type reportService interface {
	Send(string, string, string) error
}

type searchService struct {
	sessionService
	queryService
}

type Crowsnest struct {
	jobs       []job
	httpClient *http.Client
}

func New(configPath string, httpClient *http.Client) Crowsnest {
	jobs := BuildJobsFromConfig(configPath)

	return Crowsnest{jobs, httpClient}
}

func (cn *Crowsnest) ScheduleJobs() {
	for _, j := range cn.jobs {
		sessionService := session.New(j.Search.Host, j.Search.getUsername(), j.Search.getPassword(), cn.httpClient)

		query := j.NewSearch(cn.httpClient)

		searchService := searchService{sessionService, query}

		s.Every(j.Frequency).Minutes().Do(j.GetCron(searchService, report.Report{}))
	}
}

func (cn Crowsnest) StartBlocking() {
	s.StartBlocking()
}
