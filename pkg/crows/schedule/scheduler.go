package schedule

import (
	"log"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/go-co-op/gocron"
)

type Schedule struct {
	Gocron    *gocron.Scheduler
	DelayJobs int
}

func (s *Schedule) Load(list *job.List) {
	s.Gocron.Clear()

	for i, j := range list.Jobs {
		s.Gocron.Every(j.Frequency).Minutes().Tag(j.Name).Do(j.Func())

		log.Printf("⏲️ Scheduled Job %d: %s for every %d min(s)", i, j.Name, j.Frequency)

		time.Sleep(time.Duration(s.DelayJobs) * time.Second)
	}

	s.Gocron.StartAsync()
}

func (s Schedule) NextRun(job *job.Job) string {
	return s.getCronByTag(job.Name).NextRun().Format(time.RFC822)
}

func (s Schedule) LastRun(job *job.Job) string {
	return s.getCronByTag(job.Name).LastRun().Format(time.RFC822)
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
