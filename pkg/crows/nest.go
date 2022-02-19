package crows

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/go-co-op/gocron"
)

type Nest struct {
	mu     sync.Mutex
	list   job.Lister
	gocron *gocron.Scheduler
}

func NewNest(list job.Lister, goc *gocron.Scheduler) *Nest {
	goc.StartAsync()

	n := &Nest{list: list, gocron: goc}

	for _, j := range list.Jobs() {
		n.add(j.Name, j.Frequency, j.GetOffSetTime(), j.GetFunc(), true)
	}

	return n
}

func (n *Nest) HandleEvent(event job.Event) {
	go func(n *Nest, event job.Event) {
		n.mu.Lock()
		defer n.mu.Unlock()

		log.Printf("inbound event: %#v", event)

		n.list.HandleEvent(event)

		for _, j := range n.list.Jobs() {
			n.add(j.Name, j.Frequency, j.GetOffSetTime(), j.GetFunc(), true)
		}
	}(n, event)
}

func (n *Nest) add(name string, frequency int, startAt time.Time, do func(), replaceExisting bool) {
	existingJob, err := n.getCronByTag(name)
	if err == nil && existingJob.IsRunning() {
		log.Printf("Job is already running with this tag: %s, replace: %t", name, replaceExisting)
		if !replaceExisting {
			return
		}
		n.gocron.Remove(existingJob)
	}

	log.Printf("schedule %s every %d min(s) to begin at %s", name, frequency, startAt)

	n.gocron.Every(frequency).Minutes().StartAt(startAt).Tag(name).Do(do)
}

func (n *Nest) getCronByTag(tag string) (*gocron.Job, error) {
	for _, cj := range n.gocron.Jobs() {
		for _, t := range cj.Tags() {
			if tag == t {
				return cj, nil
			}
		}
	}
	return &gocron.Job{}, errors.New("no job found for tag:" + tag)
}

func (n *Nest) Jobs() []*job.Job {
	return n.list.Jobs()
}

func (n *Nest) NextRun(name string) time.Time {
	j, err := n.getCronByTag(name)
	if err != nil {
		return time.Now()
	}
	return j.NextRun()
}

func (n *Nest) LastRun(name string) time.Time {
	j, err := n.getCronByTag(name)
	if err != nil {
		return time.Now()
	}
	return j.LastRun()
}
