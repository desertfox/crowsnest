package schedule

import (
	"log"
	"net/http"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/go-co-op/gocron"
)

var httpClient *http.Client

type Schedule struct {
	gocron *gocron.Scheduler
	delay  int
}

func (s *Schedule) Load(delay int) *Schedule {
	httpClient = &http.Client{}

	s = &Schedule{
		gocron: gocron.NewScheduler(time.UTC),
		delay:  delay,
	}
	return s
}

//list needs its own struct to encapsulate un/pw/client being provided for
//constructing search service and output service.
func (s *Schedule) Run(list job.List, un, pw string) {
	for i, j := range list {
		//list.getFunc(i)
		jobFunc := j.Func(
			j.Search.Service(
				un,
				pw,
				httpClient,
			),
			j.Output.Service(),
		)

		//frequency needs to move to value attrib of job from search.
		s.gocron.Every(j.Search.Frequency).Minutes().Tag(j.Name).Do(jobFunc)

		log.Printf("⏲️ Scheduled Job %d: %s for every %d min(s)", i, j.Name, j.Search.Frequency)

		time.Sleep(time.Duration(s.delay) * time.Second)
	}

	s.gocron.StartAsync()
}

func (s *Schedule) ClearAndRun(list job.List, un, pw string, delay int) {
	log.Printf("Schedule Clearing Jobs : %v", len(s.gocron.Jobs()))

	s.gocron.Clear()

	s.Run(list, un, pw)
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
