package schedule

import (
	"log"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/go-co-op/gocron"
)

type Schedule struct {
	gocron *gocron.Scheduler
}

func NewSchedule(goc *gocron.Scheduler) *Schedule {
	return &Schedule{
		gocron: goc,
	}
}

func (s Schedule) Load(list *job.List) {
	s.gocron.Clear()

	for _, j := range list.Jobs {
		go s.scheduleJob(j)
	}

	s.gocron.StartAsync()
}

func (s Schedule) scheduleJob(j *job.Job) {
	log.Printf("schedule Job %s for every %d min(s) to begin at %s", j.Name, j.Frequency, j.GetOffSetTime())

	s.gocron.Every(j.Frequency).Minutes().StartAt(j.GetOffSetTime()).Tag(j.Name).Do(j.Func())
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
