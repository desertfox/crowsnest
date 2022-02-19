package schedule

import (
	"log"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/go-co-op/gocron"
)

type List interface {
	Jobs() []*job.Job
}

type Scheduler interface {
	Load(List)
	NextRun(string) time.Time
	LastRun(string) time.Time
}

type Schedule struct {
	gocron *gocron.Scheduler
}

func NewSchedule(goc *gocron.Scheduler) *Schedule {
	return &Schedule{
		gocron: goc,
	}
}

func (s Schedule) Load(list List) {
	s.gocron.Clear()

	for _, j := range list.Jobs() {
		go s.schedule(j.Name, j.Frequency, j.GetOffSetTime(), j.GetFunc())
	}

	s.gocron.StartAsync()
}

func (s Schedule) schedule(name string, frequency int, startAt time.Time, do func()) {
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
