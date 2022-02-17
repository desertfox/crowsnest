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
	version string = "v2.4"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	list := &job.List{}
	list.Load()

	scheduler := &schedule.Schedule{
		Gocron: gocron.NewScheduler(time.UTC),
	}
	scheduler.Load(list)

	nest := &crows.Nest{
		List:      list,
		Scheduler: scheduler,
	}
	api.New(nest).Load()
}
