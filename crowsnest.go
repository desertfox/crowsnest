package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	delayJobs  string            = os.Getenv("CROWSNEST_DELAY")
	s          *gocron.Scheduler = gocron.NewScheduler(time.UTC)
	newJobChan chan jobs.Job     = make(chan jobs.Job)
	eventChan  chan string       = make(chan string)
)

type crowsnest struct {
	jobs *jobs.JobList
}

func main() {
	log.Println("Crowsnest Startup")

	jobList, err := jobs.BuildFromConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("error building jobs from config: %v, error: %v", configPath, err.Error()))
	}

	log.Println("Crowsnest JobRunner Startup")

	cn := crowsnest{jobList}
	cn.Run(un, pw)

	go addNewJob(newJobChan, cn, un, pw)

	go handleEvent(eventChan, &cn)

	log.Println("Crowsnest Server Startup")

	server := api.NewServer(&http.ServeMux{}, newJobChan, eventChan, s)
	server.Run()
}

func (cn crowsnest) Run(un, pw string) {
	log.Println("Crowsnest ScheduleJobs")
	cn.ScheduleJobs(un, pw)

	log.Println("Crowsnest Daemon")
	cn.StartAsync()
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
			j.Search.Frequency,
			j.Search.Fields,
			j.Search.Type,
			httpClient,
		)

		searchService := jobs.SearchService{
			SessionService: sessionService,
			QueryService:   queryService,
		}

		reportService := report.Report{
			Url: j.Output.TeamsURL,
		}

		s.Every(j.Search.Frequency).Minutes().Tag(j.Name).Do(j.GetCron(searchService, reportService))

		log.Printf("Scheduled Job %d: %s for every %d min(s)", i, j.Name, j.Search.Frequency)

		if delayJobs != "" {
			delay, err := strconv.Atoi(delayJobs)
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(time.Duration(delay) * time.Second)
		}
	}
}

func (cn crowsnest) StartAsync() {
	s.StartAsync()
}

func addNewJob(newJobChan chan jobs.Job, cn crowsnest, un, pw string) {
	job := <-newJobChan

	log.Println(fmt.Sprintf("New Job recv on channel to scheduler: %#v", job))

	cn.jobs.Add(job)

	cn.jobs.WriteConfig(configPath)

	cn.Run(un, pw)
}

func handleEvent(event chan string, cn *crowsnest, un, pw string) {
	switch <-event {
	case "reloadjobs":
		log.Println("ReloadJobs event")

		jobList, err := jobs.BuildFromConfig(configPath)
		if err != nil {
			log.Fatal(fmt.Sprintf("error building jobs from config: %v, error: %v", configPath, err.Error()))
		}

		cn.jobs = jobList

		cn.Run(un, pw)
	}
}
