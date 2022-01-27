package main

import (
	"log"
	"time"

	"github.com/desertfox/crowsnest/api"
	"github.com/desertfox/crowsnest/config"
	"github.com/desertfox/crowsnest/pkg/crows"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/crowsnest/pkg/crows/schedule"
	"github.com/go-co-op/gocron"
)

const (
	version string = "v1.6"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	config := &config.Config{}
	config.Load()

	list := &job.List{
		Config: config,
	}
	list.Load()

	scheduler := &schedule.Schedule{
		DelayJobs: config.DelayJobs,
		Gocron:    gocron.NewScheduler(time.UTC),
	}
	scheduler.Load(list)

	nest := &crows.Nest{
		List:      list,
		Scheduler: scheduler,
	}
	nest.Load()

	api.New(nest).Load()
}
