package jobs

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/teams/report"
	"github.com/go-co-op/gocron"
)

type Jobs struct {
	config         *Config
	jobList        *JobList
	event          chan Event
	scheduler      *gocron.Scheduler
	httpClient     *http.Client
	loadJobChannel sync.Once
}

func Load() *Jobs {
	config := LoadConfigFromEnv()

	jobList := LoadJobList(config.Path)

	return &Jobs{
		config:     config,
		jobList:    jobList,
		event:      make(chan Event),
		scheduler:  gocron.NewScheduler(time.UTC),
		httpClient: &http.Client{},
	}
}

func (j Jobs) EventChannel() chan Event {
	return j.event
}

func (j Jobs) Scheduler() *gocron.Scheduler {
	return j.scheduler
}

func (j *Jobs) WriteConfig() {
	j.jobList.WriteConfig(j.config.Path)
}

func (j Jobs) Schedule() {
	j.loadJobChannel.Do(func() {
		go j.HandleEvent()
	})

	if len(j.scheduler.Jobs()) > 0 {
		log.Printf("🧹 Scheduler Clear Jobs : %v", len(j.scheduler.Jobs()))
		j.scheduler.Clear()
	}

	for i, job := range *j.jobList {
		jobFunc := job.Func(
			SearchService{
				SessionService: session.New(
					job.Search.Host,
					j.config.Username,
					j.config.Password,
					j.httpClient,
				),
				QueryService: search.New(
					job.Search.Host,
					job.Search.Query,
					job.Search.Streamid,
					job.Search.Frequency,
					job.Search.Fields,
					job.Search.Type,
					j.httpClient,
				),
			},
			report.Report{
				Url: job.Output.TeamsURL,
			},
		)

		j.scheduler.Every(job.Search.Frequency).Minutes().Tag(job.Name).Do(jobFunc)

		log.Printf("⏲️ Scheduled Job %d: %s for every %d min(s)", i, job.Name, job.Search.Frequency)

		time.Sleep(time.Duration(j.config.DelayJobs) * time.Second)
	}

	j.scheduler.StartAsync()
}

func (j *Jobs) HandleEvent() {
	event := <-j.event
	switch event.Action {
	case ReloadJobList:
		j.jobList = LoadJobList(j.config.Path)
	case DelJob:
		j.jobList.Del(event.Value)
		j.WriteConfig()
	case AddJob:
		j.jobList.Add(event.Job)
		j.WriteConfig()
	}
}
