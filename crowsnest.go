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
	jobs       *jobs.JobList
	newJobChan chan jobs.Job
}

func main() {
	log.Println("Crowsnest Startup")

	jobList, err := jobs.BuildFromConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("error building jobs from config: %v, error: %v", configPath, err.Error()))
	}

	log.Println("Crowsnest JobRunner Startup")

	newJobChan := make(chan jobs.Job)
	cn := crowsnest{jobList, newJobChan}
	cn.Run()

	log.Println("Crowsnest Server Startup")

	server := api.NewServer(&http.ServeMux{}, newJobChan)
	server.Run()
}

func (cn crowsnest) Run() {
	log.Println("Crowsnest ScheduleJobs")
	cn.ScheduleJobs(un, pw)

	log.Println("Crowsnest Daemon")
	cn.StartAsync()

	go func(un, pw string) {
		job := <-cn.newJobChan

		cn.jobs.Add(job)

		cn.jobs.WriteConfig(configPath)

		log.Println(fmt.Sprintf("New Job recv on channel to scheduler %#v", cn.jobs))

		cn.ScheduleJobs(un, pw)
	}(un, pw)
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
			j.Search.Type,
			time.Now().Add(-1*time.Duration(j.Frequency)*time.Minute),
			time.Now(),
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

		log.Printf("Scheduled Job %d: %s", i, j.Name)
	}
}

func (cn crowsnest) StartAsync() {
	s.StartAsync()
}
