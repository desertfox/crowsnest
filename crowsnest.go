package main

import (
	"log"
	"net/http"
	"time"

	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/jobs"
	"github.com/desertfox/crowsnest/pkg/teams/report"
	"github.com/go-co-op/gocron"
)

const (
	version   string = "v1.2"
	logPrefix string = "(┛ಠ_ಠ)┛彡┻━┻ Crowsnest "
)

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
	log.Printf("%s Startup Version: %s", logPrefix, version)

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

	log.Printf("%s Server Startup", logPrefix)

	server := api.Server{&http.ServeMux{}, jobEvent, cn.scheduler}
	server.Run()
}

func (cn crowsnest) ScheduleJobs() {
	log.Printf("%s ScheduleJobs", logPrefix)

	if len(cn.scheduler.Jobs()) > 0 {
		log.Printf("🧹 %s Clear Jobs from scheduler: %v", logPrefix, len(cn.scheduler.Jobs()))
		cn.scheduler.Clear()
	}

	for i, j := range *cn.jobs {
		cn.scheduler.Every(j.Search.Frequency).Minutes().Tag(j.Name).Do(j.GetCron(cn.createSearchService(j), cn.createReportService(j)))

		log.Printf("⏲️ Job %d: %s for every %d min(s)", i, j.Name, j.Search.Frequency)

		time.Sleep(time.Duration(cn.config.DelayJobs) * time.Second)
	}

	log.Printf("%s StartJobs", logPrefix)
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

func (cn *crowsnest) handleJobEvent(jobEvent chan jobs.Event) {
	cn.jobs.HandleEvent(<-jobEvent, cn.config.ConfigPath)
}
