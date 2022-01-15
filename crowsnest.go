package main

import (
	"log"
	"net/http"
	"time"

	"github.com/desertfox/crowsnest/pkg/api"
	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/jobs"
	"github.com/desertfox/crowsnest/pkg/teams/report"
	"github.com/go-co-op/gocron"
)

const version string = "v1.1"

type crowsnest struct {
	jobs      *jobs.JobList
	config    *config.Env
	scheduler *gocron.Scheduler
}

var (
	httpClient *http.Client    = &http.Client{}
	jobEvent   chan jobs.Event = make(chan jobs.Event)
)

func main() {
	log.Printf("Crowsnest Startup, version: %s", version)

	env := &config.Env{}
	env.GetEnv()

	jl := &jobs.JobList{}
	jl.GetConfig(env.ConfigPath)

	cn := crowsnest{
		jobs:      jl,
		config:    env,
		scheduler: gocron.NewScheduler(time.UTC),
	}

	cn.ScheduleJobs()

	go cn.handleJobEvent(jobEvent)

	api.NewServer(&http.ServeMux{}, jobEvent, cn.scheduler).Run()
}

func (cn *crowsnest) handleJobEvent(jobEvent chan jobs.Event) {
	cn.jobs.HandleEvent(<-jobEvent, cn.config.ConfigPath)
}

func (cn crowsnest) ScheduleJobs() {
	log.Println("Crowsnest ScheduleJobs")

	log.Printf("Crowsnest clearing jobs from scheduler: %v", len(cn.scheduler.Jobs()))
	cn.scheduler.Clear()

	for i, j := range *cn.jobs {
		cn.scheduler.Every(j.Search.Frequency).Minutes().Tag(j.Name).Do(j.GetCron(cn.createSearchService(j), cn.createReportService(j)))

		log.Printf("Scheduled Job %d: %s for every %d min(s)", i, j.Name, j.Search.Frequency)

		time.Sleep(time.Duration(cn.config.DelayJobs) * time.Second)
	}

	log.Println("Crowsnest StartJobs")
	cn.scheduler.StartAsync()
}

func (cn crowsnest) createSearchService(j jobs.Job) jobs.SearchService {
	sessionService := session.New(
		j.Search.Host,
		cn.config.Username,
		cn.config.Password,
		httpClient,
	)

	queryService := search.New(
		j.Search.Host,
		j.Search.Query,
		j.Search.Streamid,
		j.Search.Frequency,
		j.Search.Fields,
		j.Search.Type,
		httpClient,
	)

	return jobs.SearchService{
		SessionService: sessionService,
		QueryService:   queryService,
	}
}

func (cn crowsnest) createReportService(j jobs.Job) report.Report {
	return report.Report{
		Url: j.Output.TeamsURL,
	}
}
