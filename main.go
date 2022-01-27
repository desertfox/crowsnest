package main

import (
	"log"
	"os"
	"time"

	"github.com/desertfox/crowsnest/api"
	"github.com/desertfox/crowsnest/config"
	"github.com/desertfox/crowsnest/pkg/crows"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/crowsnest/pkg/crows/schedule"
	"github.com/go-co-op/gocron"
)

const (
	version string = "v1.7"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	config := config.Load()

	var list *job.List
	if os.Getenv("CROWSNEST_DEBUG") == "1" {
		list = &job.List{}
	} else {
		list = &job.List{
			Config: config,
		}
		list.Load()
	}

	scheduler := &schedule.Schedule{
		DelayJobs: config.DelayJobs,
		Gocron:    gocron.NewScheduler(time.UTC),
	}
	scheduler.Load(list)

	nest := &crows.Nest{
		List:          list,
		Scheduler:     scheduler,
		EventCallback: make(chan crows.Event),
	}
	nest.Load()

	api.New(nest).Load()
}
