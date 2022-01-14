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

const version string = "v1.0"

var (
	httpClient *http.Client      = &http.Client{}
	un         string            = os.Getenv("CROWSNEST_USERNAME")
	pw         string            = os.Getenv("CROWSNEST_PASSWORD")
	configPath string            = os.Getenv("CROWSNEST_CONFIG")
	delayJobs  string            = os.Getenv("CROWSNEST_DELAY")
	s          *gocron.Scheduler = gocron.NewScheduler(time.UTC)
	jobChan    chan jobs.Job     = make(chan jobs.Job)
	eventChan  chan jobs.Event   = make(chan jobs.Event)
)

type crowsnest struct {
	jobs *jobs.JobList
}

func main() {
	log.Printf("Crowsnest Startup, version:%s", version)

	jobList, err := jobs.BuildFromConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("error building jobs from config: %v, error: %v", configPath, err.Error()))
	}

	cn := crowsnest{jobList}
	cn.Run(un, pw)

	go handleAddNewJob(jobChan, &cn, un, pw)

	go handleEvent(eventChan, &cn, un, pw)

	server := api.NewServer(&http.ServeMux{}, jobChan, eventChan, s)
	server.Run()
}

func (cn crowsnest) Run(un, pw string) {
	log.Println("Crowsnest Run")

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

func handleAddNewJob(newJobChan chan jobs.Job, cn *crowsnest, un, pw string) {
	job := <-newJobChan

	log.Println(fmt.Sprintf("New Job recv on channel to scheduler: %#v", job))

	cn.jobs.Add(job)

	cn.jobs.WriteConfig(configPath)

	cn.Run(un, pw)
}

func handleEvent(event chan jobs.Event, cn *crowsnest, un, pw string) {
	e := <-event

	switch e.Action {
	case "reloadjobs":
		log.Println("ReloadJobs event")

		jobList, err := jobs.BuildFromConfig(configPath)
		if err != nil {
			log.Fatal(fmt.Sprintf("error building jobs from config: %v, error: %v", configPath, err.Error()))
		}

		cn.jobs = jobList

		cn.Run(un, pw)
	case "DEL_TAG":
		tag := e.Value

		newJobList := cn.jobs.Del(tag)

		cn.jobs = &newJobList

		cn.jobs.WriteConfig(configPath)

		cn.Run(un, pw)
	}
}
