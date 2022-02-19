package schedule

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler interface {
	Add(string, int, time.Time, func(), bool)
	NextRun(string) time.Time
	LastRun(string) time.Time
}

type Schedule struct {
	gocron *gocron.Scheduler
}

func NewSchedule(goc *gocron.Scheduler) *Schedule {
	goc.StartAsync()

	return &Schedule{
		gocron: goc,
	}
}

func (s Schedule) Add(name string, frequency int, startAt time.Time, do func(), replaceExisting bool) {
	existingJob := s.getCronByTag(name)

	if existingJob.IsRunning() {
		log.Printf("Job is already running with this tag: %s, replace: %t", name, replaceExisting)
		if !replaceExisting {
			return
		}
		s.gocron.Remove(existingJob)
	}

	log.Printf("schedule %s every %d min(s) to begin at %s", name, frequency, startAt)

	s.gocron.Every(frequency).Minutes().StartAt(startAt).Tag(name).Do(do)
}

func (s Schedule) NextRun(name string) time.Time {
	return s.getCronByTag(name).NextRun()
}

func (s Schedule) LastRun(name string) time.Time {
	return s.getCronByTag(name).LastRun()
}

func (s Schedule) getCronByTag(tag string) *gocron.Job {
	for _, cj := range s.gocron.Jobs() {
		for _, t := range cj.Tags() {
			if tag == t {
				return cj
			}
		}
	}
	return &gocron.Job{}
}
