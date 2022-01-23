package main

import (
	"log"

	"github.com/desertfox/crowsnest/pkg/api"
	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/crows"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/crowsnest/pkg/crows/schedule"
)

const (
	version string = "v1.6"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	config := &config.Config{}
	config.Load()

	list := &job.List{}
	list.Load(config)

	scheduler := &schedule.Schedule{}
	scheduler.Load(config.DelayJobs)

	nest := &crows.Nest{}
	nest.Load(config, list, scheduler)
	nest.Run()

	api.New(nest).Run()
}
