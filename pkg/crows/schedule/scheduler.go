package schedule

import (
	"log"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/go-co-op/gocron"
)

type Schedule struct {
	Gocron *gocron.Scheduler
}

func (s Schedule) Load(list *job.List) {
	s.Gocron.Clear()

	for _, j := range list.Jobs {
		go s.scheduleJob(j)
	}

	s.Gocron.StartAsync()
}

func (s Schedule) scheduleJob(j *job.Job) {
	log.Printf("⏲️ Schedule Job %s for every %d min(s) to begin at %s", j.Name, j.Frequency, j.GetOffSetTime())

	s.Gocron.Every(j.Frequency).Minutes().StartAt(j.GetOffSetTime()).Tag(j.Name).Do(j.Func())
}

func (s Schedule) NextRun(job *job.Job) time.Time {
	return s.getCronByTag(job.Name).NextRun()
}

func (s Schedule) LastRun(job *job.Job) time.Time {
	return s.getCronByTag(job.Name).LastRun()
}

func (s Schedule) getCronByTag(tag string) *gocron.Job {
	for _, cj := range s.Gocron.Jobs() {
		for _, t := range cj.Tags() {
			if tag == t {
				return cj
			}
		}
	}
	return &gocron.Job{}
}
