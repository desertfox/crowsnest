package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/jobs"
	"github.com/desertfox/crowsnest/pkg/teams/report"
	"github.com/fatih/color"
	"github.com/go-co-op/gocron"
)

var (
	httpClient *http.Client      = &http.Client{}
	un         string            = os.Getenv("CROWSNEST_USERNAME")
	pw         string            = os.Getenv("CROWSNEST_PASSWORD")
	configPath string            = os.Getenv("CROWSNEST_CONFIG")
	s          *gocron.Scheduler = gocron.NewScheduler(time.UTC)
)

type crowsnest struct {
	jobs       []jobs.Job
	httpClient *http.Client
}

func main() {
	color.Yellow("Crowsnest Startup")

	jobs := jobs.BuildFromConfig(configPath)

	for i, job := range jobs {
		color.Yellow(fmt.Sprintf("Loaded Job %d: %s", i, job.Name))
	}

	cn := crowsnest{jobs, httpClient}

	color.Yellow("Crowsnest ScheduleJobs")

	cn.ScheduleJobs(un, pw)

	color.Green("Crowsnest Daemon...")

	cn.StartBlocking()
}

func (cn *crowsnest) ScheduleJobs(un, pw string) {
	for i, j := range cn.jobs {
		sessionService := session.New(
			j.Search.Host,
			un,
			pw,
			httpClient,
		)

		queryService := search.New(
			j.Search.Host,
			j.Search.Query,
			j.Search.Streamid,
			j.Frequency,
			j.Search.Fields,
			httpClient,
		)

		searchService := jobs.SearchService{
			SessionService: sessionService,
			QueryService:   queryService,
		}

		reportService := report.Report{
			Url: j.TeamsURL,
		}

		s.Every(j.Frequency).Minutes().Do(j.GetCron(searchService, reportService))

		color.Green(fmt.Sprintf("Scheduled Job %d: %s", i, j.Name))
	}
}

func (cn crowsnest) StartBlocking() {
	s.StartBlocking()
}
