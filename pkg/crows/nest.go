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

type Event struct {
	Action int
	Value  string
	Job    *job.Job
}

type Scheduler interface {
	Load(*job.List)
	NextRun(string) time.Time
	LastRun(string) time.Time
}

type Nest struct {
	mu        sync.Mutex
	list      *job.List
	scheduler *schedule.Schedule
}

func NewNest(list *job.List, scheduler *schedule.Schedule) *Nest {
	return &Nest{
		list:      list,
		scheduler: scheduler,
	}
}

//All the nest methods bellow are used to expose schedule and job state to API
func (n *Nest) HandleEvent(event Event) {
	go func(n *Nest, event Event) {
		n.mu.Lock()
		defer n.mu.Unlock()

		log.Printf("inbound event: %#v", event)

		switch event.Action {
		case Reload:
			n.list.Clear()
			n.list.Load()
		case Del:
			n.list.Del(event.Job)
		case Add:
			n.list.Add(event.Job)
		}

		n.list.Save()

		n.scheduler.Load(n.list)
	}(n, event)
}

func (n *Nest) Jobs() []*job.Job {
	return n.list.Jobs
}

func (n *Nest) NextRun(name string) time.Time {
	return n.scheduler.NextRun(name)
}

func (n *Nest) LastRun(name string) time.Time {
	return n.scheduler.LastRun(name)
}
