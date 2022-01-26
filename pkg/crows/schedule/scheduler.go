package schedule

import (
	"log"
	"time"

	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/go-co-op/gocron"
)

type Schedule struct {
	Gocron *gocron.Scheduler
	Config *config.Config
}

func (s *Schedule) Load(list *job.List) {
	for i, j := range list.Jobs {
		s.Gocron.Every(j.Frequency).Minutes().Tag(j.Name).Do(j.Func())

		log.Printf("⏲️ Scheduled Job %d: %s for every %d min(s)", i, j.Name, j.Frequency)

		time.Sleep(time.Duration(s.Config.DelayJobs) * time.Second)
	}

	s.Gocron.StartAsync()
}

func (s *Schedule) ClearAndLoad(list *job.List) {
	log.Printf("Schedule Clearing Jobs : %v", len(s.Gocron.Jobs()))

	s.Gocron.Clear()

	s.Load(list)
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
