package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog/cron"
	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/teams"
	"github.com/go-co-op/gocron"
)

func main() {

	httpClient := &http.Client{}

	var host string = os.Getenv("CROWSNEST_HOST")

	sessionService, err := session.NewSession(
		host,
		os.Getenv("CROWSNEST_USERNAME"),
		os.Getenv("CROWSNEST_PASSWORD"),
		httpClient,
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	jobService := cron.BuildFromConfig(os.Getenv("CROWSNEST_CONFIG"))

	s := gocron.NewScheduler(time.UTC)

	for _, job := range *jobService {
		outputService := teams.BuildClient(job.TeamsURL)

		query := search.NewQuery(host, job.Name, job.Option.Query, job.Option.Streamid, job.Frequency, job.Option.Fields)

		s.Every(job.Frequency).Minutes().Do(job.GetFunc(sessionService, outputService, query))
	}

	s.StartBlocking()
}
