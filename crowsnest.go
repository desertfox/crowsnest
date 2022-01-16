package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/jobs"
	"github.com/desertfox/crowsnest/pkg/teams/report"
	"github.com/go-co-op/gocron"
)

type crowsnest struct {
	jobs        *jobs.JobList
	config      *config.Config
	scheduler   *gocron.Scheduler
	event       chan jobs.Event
	httpClient  *http.Client
	loadChannel sync.Once
}

func (cn crowsnest) ScheduleJobs() {
	crowLog("ScheduleJobs")

	cn.loadChannel.Do(func() {
		go cn.handleJobEvent(cn.event)
	})

	if len(cn.scheduler.Jobs()) > 0 {
		crowLog(fmt.Sprintf("🧹 Scheduler Clear Jobs : %v", len(cn.scheduler.Jobs())))
		cn.scheduler.Clear()
	}

	for i, j := range *cn.jobs {
		jobFunc := j.GetCron(
			jobs.SearchService{
				SessionService: session.New(
					j.Search.Host,
					cn.config.Username,
					cn.config.Password,
					cn.httpClient,
				),
				QueryService: search.New(
					j.Search.Host,
					j.Search.Query,
					j.Search.Streamid,
					j.Search.Frequency,
					j.Search.Fields,
					j.Search.Type,
					cn.httpClient,
				),
			},
			report.Report{
				Url: j.Output.TeamsURL,
			},
		)

		cn.scheduler.Every(j.Search.Frequency).Minutes().Tag(j.Name).Do(jobFunc)

		log.Printf("⏲️ Scheduled Job %d: %s for every %d min(s)", i, j.Name, j.Search.Frequency)

		time.Sleep(time.Duration(cn.config.DelayJobs) * time.Second)
	}

	cn.scheduler.StartAsync()
}

func (cn *crowsnest) handleJobEvent(jobEvent chan jobs.Event) {
	event := <-jobEvent

	crowLog(fmt.Sprintf("handleJobEvent event %v", event))

	cn.jobs.HandleEvent(event, cn.config.Path)
}
