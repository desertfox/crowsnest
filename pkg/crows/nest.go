package crows

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/go-co-op/gocron"
)

var (
	wg sync.WaitGroup
	mu sync.Mutex
)

type Nest struct {
	list   job.Lister
	gocron *gocron.Scheduler
}

func NewNest(list job.Lister, goc *gocron.Scheduler) *Nest {
	n := &Nest{list: list, gocron: goc}

	n.gocron.StartAsync()

	n.schedule()

	return n
}

func (n *Nest) HandleEvent(event job.Event) {
	go func(event job.Event) {
		log.Printf("inbound event: %#v", event)

		n.list.HandleEvent(event)

		n.schedule()
	}(event)
}

func (n *Nest) schedule() {
	if len(n.list.Jobs()) == 0 {
		return
	}

	for _, j := range n.list.Jobs() {
		wg.Add(1)
		go func(name string, frequency int, startAt time.Time, f func()) {
			defer wg.Done()
			n.add(name, frequency, startAt, f, true)
		}(j.Name, j.Frequency, j.GetOffSetTime(), j.GetFunc())
	}
	wg.Wait()
}

func (n *Nest) add(name string, frequency int, startAt time.Time, do func(), replaceExisting bool) {
	existingJob, err := n.getCronByTag(name)
	if err == nil && existingJob.IsRunning() {
		log.Printf("Job is already running with this tag: %s, replace: %t", name, replaceExisting)
		if !replaceExisting {
			return
		}
		mu.Lock()
		n.gocron.Remove(existingJob)
		mu.Unlock()
	}

	log.Printf("schedule %s every %d min(s) to begin at %s", name, frequency, startAt)

	mu.Lock()
	n.gocron.Every(frequency).Minutes().StartAt(startAt).Tag(name).Do(do)
	mu.Unlock()
}

func (n *Nest) getCronByTag(tag string) (*gocron.Job, error) {
	mu.Lock()
	defer mu.Unlock()
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
