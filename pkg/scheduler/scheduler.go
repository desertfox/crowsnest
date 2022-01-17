package scheduler

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/job"
	"github.com/desertfox/crowsnest/pkg/joblist"
	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	config         *config.Config
	jobList        joblist.JobList
	event          chan job.Event
	scheduler      *gocron.Scheduler
	httpClient     *http.Client
	loadJobChannel sync.Once
}

func Load() *Scheduler {
	config := config.LoadConfigFromEnv()

	jobList := joblist.JobList{}

	return &Scheduler{
		config:     config,
		jobList:    jobList.Load(config.Path),
		event:      make(chan job.Event),
		scheduler:  gocron.NewScheduler(time.UTC),
		httpClient: &http.Client{},
	}
}

func (js Scheduler) EventChannel() chan job.Event {
	return js.event
}

func (js Scheduler) Jobs() joblist.JobList {
	return js.jobList
}

func (js Scheduler) SJobs() []*gocron.Job {
	return js.scheduler.Jobs()
}

func (js *Scheduler) Save() {
	js.jobList.Save(js.config.Path)
}

func (js *Scheduler) Schedule() *Scheduler {
	js.loadJobChannel.Do(func() {
		go js.HandleEvent()
	})

	if len(js.scheduler.Jobs()) > 0 {
		log.Printf("Scheduler Clearing Jobs : %v", len(js.scheduler.Jobs()))
		js.scheduler.Clear()
	}

	for i, j := range js.jobList {
		jobFunc := j.Func(
			j.Search.SearchService(
				js.config.Username,
				js.config.Password,
				js.httpClient,
			),
			j.Output.ReportService(),
		)

		js.scheduler.Every(j.Search.Frequency).Minutes().Tag(j.Name).Do(jobFunc)

		log.Printf("⏲️ Scheduled Job %d: %s for every %d min(s)", i, j.Name, j.Search.Frequency)

		time.Sleep(time.Duration(js.config.DelayJobs) * time.Second)
	}

	js.scheduler.StartAsync()

	return js
}

func (js *Scheduler) HandleEvent() {
	event := <-js.event
	switch event.Action {
	case job.Reload:
		js.jobList = joblist.JobList{}
		js.jobList.Load(js.config.Path)
	case job.Del:
		js.jobList.Del(event.Value)
		js.Save()
	case job.Add:
		js.jobList.Add(event.Job)
		js.Save()
	}
	js.Schedule()
}
