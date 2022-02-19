package crows

import (
	"log"
	"sync"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/crowsnest/pkg/crows/schedule"
)

type Nest struct {
	mu       sync.Mutex
	list     job.Lister
	schedule schedule.Scheduler
}

func NewNest(list job.Lister, scheduler schedule.Scheduler) *Nest {
	return &Nest{
		list:     list,
		schedule: scheduler,
	}
}

func (n *Nest) HandleEvent(event job.Event) {
	go func(n *Nest, event job.Event) {
		n.mu.Lock()
		defer n.mu.Unlock()

		log.Printf("inbound event: %#v", event)

		n.list.HandleEvent(event)

		for _, j := range n.list.Jobs() {
			n.schedule.Add(j.Name, j.Frequency, j.GetOffSetTime(), j.GetFunc(), true)
		}
	}(n, event)
}

func (n *Nest) Jobs() []*job.Job {
	return n.list.Jobs()
}

func (n *Nest) NextRun(name string) time.Time {
	return n.schedule.NextRun(name)
}

func (n *Nest) LastRun(name string) time.Time {
	return n.schedule.LastRun(name)
}
