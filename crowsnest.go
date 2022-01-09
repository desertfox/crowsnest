package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/desertfox/crowsnest/pkg/api"
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
	configPath string            = "config.yaml-example" //os.Getenv("CROWSNEST_CONFIG")
	s          *gocron.Scheduler = gocron.NewScheduler(time.UTC)
)

type crowsnest struct {
	jobs       *jobs.JobList
	newJobChan chan jobs.Job
}

func main() {

	color.Yellow("Crowsnest Startup")

	jobList, err := jobs.BuildFromConfig(configPath)
	if err != nil {
		panic(err.Error())
	}

	color.Yellow("Crowsnest JobRunner Startup")

	newJobChan := make(chan jobs.Job)
	cn := crowsnest{jobList, newJobChan}
	cn.Run()

	color.Yellow("Crowsnest Server Startup")

	server := api.NewServer(&http.ServeMux{}, configPath, jobList, newJobChan)
	server.Run()
}

func (cn crowsnest) Run() {
	color.Yellow("Crowsnest ScheduleJobs")
	cn.ScheduleJobs(un, pw)

	color.Green("Crowsnest Daemon...")
	cn.StartAsync()

	go func() {
		job := <-cn.newJobChan

		cn.jobs.Add(job)

		cn.jobs.WriteConfig(configPath)

		color.Yellow("Crowsnest New Job recv on channel to scheduler")
		log.Println(fmt.Sprintf("%#v", cn.jobs))

		cn.ScheduleJobs(un, pw)
	}()
}

func (cn crowsnest) ScheduleJobs(un, pw string) {
	s.Clear()

	for i, j := range *cn.jobs {
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

func (cn crowsnest) StartAsync() {
	s.StartAsync()
}
