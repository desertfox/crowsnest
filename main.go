package main

import (
	"os"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/teams"
	"github.com/go-co-op/gocron"
)

func main() {

	sessionService := session.NewLoginRequest(
		os.Getenv("CROWSNEST_HOST"),
		os.Getenv("CROWSNEST_USERNAME"),
		os.Getenv("CROWSNEST_PASSWORD"),
	)

	jobService := graylog.NewClient(sessionService).BuildJobsFromConfig(os.Getenv("CROWSNEST_CONFIG"))

	s := gocron.NewScheduler(time.UTC)

	for _, job := range jobService.Jobs {
		outputService := teams.BuildClient(job.TeamsURL)

		s.Every(job.Frequency).Minutes().Do(job.GetFunc(jobService, outputService))
	}

	s.StartBlocking()
}
