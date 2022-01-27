package crows

import (
	"log"
	"sync"
	"time"

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
func (n *Nest) Load() {
	loadEventCallbackOnce.Do(func() {
		go func(n *Nest) {
			event := <-n.EventCallback

			log.Printf("event %#v", event)

			switch event.Action {
			case Reload:
				n.List.Clear()
				n.List.Load()
			case Del:
				n.List.Del(event.Job)
			case Add:
				n.List.Add(event.Job)
			}

			n.List.Save()

			n.Scheduler.Load(n.List)
		}(n)
	})
}

func (n Nest) Jobs() []*job.Job {
	return n.List.Jobs
}

func (n Nest) NextRun(job *job.Job) time.Time {
	return n.Scheduler.NextRun(job)
}

func (n Nest) LastRun(job *job.Job) time.Time {
	return n.Scheduler.LastRun(job)
}
