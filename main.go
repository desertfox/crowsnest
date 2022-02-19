package main

import (
	"log"
	"time"

	"github.com/desertfox/crowsnest/api"
	"github.com/desertfox/crowsnest/pkg/crows"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/crowsnest/pkg/crows/schedule"
	"github.com/go-co-op/gocron"
)

const (
	version string = "v2.7"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	l := job.NewList()

	s := schedule.NewSchedule(gocron.NewScheduler(time.UTC))

	for _, j := range l.Jobs() {
		s.Add(j.Name, j.Frequency, j.GetOffSetTime(), j.GetFunc(), true)
	}

	nest := crows.NewNest(l, s)

	api.New(nest).Load()
}
