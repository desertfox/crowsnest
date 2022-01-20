package crows

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/config"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	config         *config.Config
	list           job.List
	event          chan job.Event
	scheduler      *gocron.Scheduler
	httpClient     *http.Client
	loadJobChannel sync.Once
}

func Load() *Scheduler {
	config := config.LoadFromEnv()

	list := job.List{}

	return &Scheduler{
		config:     config,
		list:       list.Load(config.Path),
		event:      make(chan job.Event),
		scheduler:  gocron.NewScheduler(time.UTC),
		httpClient: &http.Client{},
	}
}

func (s Scheduler) EventChannel() chan job.Event {
	return s.event
}

func (s Scheduler) Jobs() job.List {
	return s.list
}

func (s Scheduler) SJobs() []*gocron.Job {
	return s.scheduler.Jobs()
}

func (s *Scheduler) Savelist() {
	s.list.Save(s.config.Path)
}

func (s *Scheduler) Schedule() *Scheduler {
	s.loadJobChannel.Do(func() {
		go s.Event()
	})

	if len(s.scheduler.Jobs()) > 0 {
		log.Printf("Scheduler Clearing Jobs : %v", len(s.scheduler.Jobs()))
		s.scheduler.Clear()
	}

	for i, j := range s.list {
		jobFunc := j.Func(
			j.Search.Service(
				s.config.Username,
				s.config.Password,
				s.httpClient,
			),
			j.Output.Service(),
		)

		s.scheduler.Every(j.Search.Frequency).Minutes().Tag(j.Name).Do(jobFunc)

		log.Printf("⏲️ Scheduled Job %d: %s for every %d min(s)", i, j.Name, j.Search.Frequency)

		time.Sleep(time.Duration(s.config.DelayJobs) * time.Second)
	}

	s.scheduler.StartAsync()

	return s
}

func (js *Scheduler) Event() {
	js.list.HandleEvent(<-js.event)

	js.Savelist()

	js.Schedule()
}
