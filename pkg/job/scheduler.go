package job

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/desertfox/crowsnest/pkg/job/config"
	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	config         *config.Config
	list           List
	event          chan Event
	scheduler      *gocron.Scheduler
	httpClient     *http.Client
	loadJobChannel sync.Once
}

func Load() *Scheduler {
	config := config.LoadConfigFromEnv()

	list := List{}

	return &Scheduler{
		config:     config,
		list:       list.Load(config.Path),
		event:      make(chan Event),
		scheduler:  gocron.NewScheduler(time.UTC),
		httpClient: &http.Client{},
	}
}

func (js Scheduler) EventChannel() chan Event {
	return js.event
}

func (js Scheduler) Jobs() List {
	return js.list
}

func (js Scheduler) SJobs() []*gocron.Job {
	return js.scheduler.Jobs()
}

func (js *Scheduler) Savelist() {
	js.list.Save(js.config.Path)
}

func (js *Scheduler) Schedule() *Scheduler {
	js.loadJobChannel.Do(func() {
		go js.HandleEvent()
	})

	if len(js.scheduler.Jobs()) > 0 {
		log.Printf("Scheduler Clearing Jobs : %v", len(js.scheduler.Jobs()))
		js.scheduler.Clear()
	}

	for i, j := range js.list {
		jobFunc := j.Func(
			j.Search.Service(
				js.config.Username,
				js.config.Password,
				js.httpClient,
			),
			j.Output.Service(),
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

	log.Printf("Event %#v", event)

	switch event.Action {
	case Reload:
		js.list = List{}
	case Del:
		js.list.Del(event.Value)
	case Add:
		js.list.Add(event.Job)
	}

	js.Savelist()

	js.Schedule()
}
