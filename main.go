package main

import (
	"log"
	"net/http"
	"time"

	"github.com/desertfox/crowsnest/pkg/api"
	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/jobs"
	"github.com/go-co-op/gocron"
)

const (
	version   string = "v1.2"
	logPrefix string = "Crowsnest "
)

func main() {
	crowLog("Startup Version " + version)

	conf := config.New()

	jobList := jobs.Load(conf.Path)

	cn := crowsnest{
		config:     conf,
		jobs:       jobList,
		scheduler:  gocron.NewScheduler(time.UTC),
		event:      make(chan jobs.Event),
		httpClient: &http.Client{},
	}
	cn.ScheduleJobs()

	crowLog("Server Startup")

	server := api.Server{
		Mux:       &http.ServeMux{},
		Event:     cn.event,
		Scheduler: cn.scheduler,
	}
	server.Run()
}

func crowLog(s string) {
	log.Printf("%s : %s", logPrefix, s)
}
