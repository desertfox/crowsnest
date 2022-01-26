package crows

import (
	"sync"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/crowsnest/pkg/crows/schedule"
)

var (
	loadEventCallbackOnce sync.Once
)

type action int

const (
	Add action = iota
	Del
	Reload
)

type Event struct {
	Action action
	Value  string
	Job    *job.Job
}

type Nest struct {
	List          *job.List
	Scheduler     *schedule.Schedule
	EventCallback chan Event
}

//All the nest methods bellow are used to expose schedule and job state to API
func (n Nest) Load() {
	loadEventCallback(n)
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

func loadEventCallback(n Nest) {
	loadEventCallbackOnce.Do(func() {
		go func() {
			n := n
			event := <-n.EventCallback

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

			n.Scheduler.ClearAndLoad(n.List)
		}()
	})
}
