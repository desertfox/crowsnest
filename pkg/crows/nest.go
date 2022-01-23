package crows

import (
	"sync"
	"time"

	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/crowsnest/pkg/crows/schedule"
)

var (
	loadEventCallback sync.Once
)

type Nest struct {
	config        *config.Config
	list          job.List
	scheduler     *schedule.Schedule
	eventCallback chan job.Event
}

func Load(config *config.Config) *Nest {
	list := job.List{}
	list.Load(config.Path)

	scheduler := &schedule.Schedule{}

	return &Nest{
		config:        config,
		list:          list,
		scheduler:     scheduler.Load(config.DelayJobs),
		eventCallback: make(chan job.Event),
	}
}

func (n Nest) Run() {
	n.loadEventCallback()

	n.scheduler.Run(n.list, n.config.Username, n.config.Password)
}

/*
All the nest methods bellow are used to expose schedule and job state to API
*/

func (n Nest) Jobs() job.List {
	return n.list
}

func (n Nest) NextRun(job *job.Job) time.Time {
	return n.scheduler.NextRun(job)
}

func (n Nest) LastRun(job *job.Job) time.Time {
	return n.scheduler.LastRun(job)
}

func (n Nest) EventCallback() chan job.Event {
	return n.eventCallback
}

func (n Nest) loadEventCallback() {
	loadEventCallback.Do(func() {
		go func() {
			n := n
			event := <-n.eventCallback

			switch event.Action {
			case job.Reload:
				n.list = job.List{}
			case job.Del:
				n.list.Del(event.Job)
			case job.Add:
				n.list.Add(event.Job)
			}

			n.list.Save(n.config.Path)

			n.scheduler.ClearAndRun(n.list, n.config.Username, n.config.Password, n.config.DelayJobs)
		}()
	})
}
