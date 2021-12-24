package main

import (
	"net/http"
	"os"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog"
	"github.com/desertfox/crowsnest/pkg/graylog/cron"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/teams"
	"github.com/go-co-op/gocron"
)

func main() {

	httpClient := &http.Client{}

	sessionService := session.NewLoginRequest(
		os.Getenv("CROWSNEST_HOST"),
		os.Getenv("CROWSNEST_USERNAME"),
		os.Getenv("CROWSNEST_PASSWORD"),
		httpClient,
	)

	jobService := cron.BuildFromConfig(os.Getenv("CROWSNEST_CONFIG"))

	graylogClient := graylog.NewClient(sessionService)

	s := gocron.NewScheduler(time.UTC)

	for _, job := range *jobService {
		outputService := teams.BuildClient(job.TeamsURL)

		s.Every(job.Frequency).Minutes().Do(job.GetFunc(graylogClient, outputService))
	}

	s.StartBlocking()
}
