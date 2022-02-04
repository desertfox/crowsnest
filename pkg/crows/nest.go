package crows

import (
	"log"
	"sync"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/crowsnest/pkg/crows/schedule"
)

const (
	Add int = iota
	Del
	Reload
)

var (
	loadEventCallbackOnce sync.Once
)

type Event struct {
	Action int
	Value  string
	Job    *job.Job
}

type Nest struct {
	List          *job.List
	Scheduler     *schedule.Schedule
	EventCallback chan Event
}

//All the nest methods bellow are used to expose schedule and job state to API
func (n *Nest) EventChannel() chan Event {
	go func(n *Nest) {
		event := <-n.EventCallback

		log.Printf("inbound event:%#v", event)

		switch event.Action {
		case Reload:
			n.List.Clear()
			n.List.Load()
		case Del:
			n.List.Del(event.Job)
		case Add:
			n.List.Add(event.Job)
		}

		log.Printf("save list:%#v", n.List)

		n.List.Save()

		n.Scheduler.Load(n.List)
	}(n)

	return n.EventCallback
}

func (n Nest) Jobs() []*job.Job {
	return n.List.Jobs
}

func (n Nest) NextRun(job *job.Job) string {
	return n.Scheduler.NextRun(job)
}

func (n Nest) LastRun(job *job.Job) string {
	return n.Scheduler.LastRun(job)
}
