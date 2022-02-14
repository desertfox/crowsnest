package crows

import (
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

type Nest struct {
	mu        sync.Mutex
	List      *job.List
	Scheduler *schedule.Schedule
}

//All the nest methods bellow are used to expose schedule and job state to API
func (n *Nest) HandleEvent(event Event) {
	go func(n *Nest, event Event) {
		n.mu.Lock()
		defer n.mu.Unlock()

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
	}(n, event)
}

func (n *Nest) Jobs() []*job.Job {
	return n.List.Jobs
}

func (n *Nest) NextRun(job *job.Job) time.Time {
	return n.Scheduler.NextRun(job)
}

func (n *Nest) LastRun(job *job.Job) time.Time {
	return n.Scheduler.LastRun(job)
}
