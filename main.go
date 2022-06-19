package main

import (
	"log"
	"time"

	"github.com/desertfox/crowsnest/api"
	"github.com/desertfox/crowsnest/pkg/crows"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/go-co-op/gocron"
)

const (
	version string = "v3.0"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	nest := crows.NewNest(job.NewList(), gocron.NewScheduler(time.UTC))

	api.New(nest).Load()
}
