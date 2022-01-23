package schedule

import (
	"log"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/go-co-op/gocron"
)

type Schedule struct {
	gocron *gocron.Scheduler
	delay  int
}

func (s *Schedule) Load(delay int) *Schedule {
	s = &Schedule{
		gocron: gocron.NewScheduler(time.UTC),
		delay:  delay,
	}
	return s
}

func (s *Schedule) Run(list *job.List) {
	for i, j := range list.Jobs {
		s.gocron.Every(j.Frequency).Minutes().Tag(j.Name).Do(j.Func())

		log.Printf("⏲️ Scheduled Job %d: %s for every %d min(s)", i, j.Name, j.Frequency)

		time.Sleep(time.Duration(s.delay) * time.Second)
	}

	s.gocron.StartAsync()
}

func (s *Schedule) ClearAndRun(list *job.List) {
	log.Printf("Schedule Clearing Jobs : %v", len(s.gocron.Jobs()))

	s.gocron.Clear()

	s.Run(list)
}

func (s Schedule) NextRun(job *job.Job) time.Time {
	return s.getCronByTag(job.Name).NextRun()
}

func (s Schedule) LastRun(job *job.Job) time.Time {
	return s.getCronByTag(job.Name).LastRun()
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
