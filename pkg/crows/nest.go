package crows

import (
	"time"

	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/crowsnest/pkg/crows/schedule"
)

type Nest struct {
	config        *config.Config
	list          *job.List
	scheduler     *schedule.Schedule
	eventCallback chan Event
}

func (n *Nest) Load(config *config.Config, scheduler *schedule.Schedule, list *job.List) *Nest {
	return &Nest{
		config:        config,
		list:          list,
		scheduler:     scheduler,
		eventCallback: make(chan Event),
	}
}

func (n Nest) Run() {
	loadEventCallback(n)

	n.scheduler.Run(n.list)
}

/*
All the nest methods bellow are used to expose schedule and job state to API
*/

func (n Nest) Jobs() []*job.Job {
	return n.list.Jobs
}

func (n Nest) NextRun(job *job.Job) time.Time {
	return n.scheduler.NextRun(job)
}

func (n Nest) LastRun(job *job.Job) time.Time {
	return n.scheduler.LastRun(job)
}

func (n Nest) EventCallback() chan Event {
	return n.eventCallback
}
