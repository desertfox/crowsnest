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
	version string = "v2.5"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	list := job.NewList()
	list.Load()

	scheduler := schedule.NewSchedule(gocron.NewScheduler(time.UTC))
	scheduler.Load(list)

	nest := crows.NewNest(list, scheduler)

	api.New(nest).Load()
}
