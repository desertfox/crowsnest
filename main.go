package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog"
	"github.com/desertfox/crowsnest/pkg/graylog/cron"
	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/report"
	"github.com/go-co-op/gocron"
)

var (
	httpClient *http.Client      = &http.Client{}
	host       string            = os.Getenv("CROWSNEST_HOST")
	username   string            = os.Getenv("CROWSNEST_USERNAME")
	password   string            = os.Getenv("CROWSNEST_PASSWORD")
	configPath string            = os.Getenv("CROWSNEST_CONFIG")
	s          *gocron.Scheduler = gocron.NewScheduler(time.UTC)
)

func main() {
	sessionService, err := session.New(host, username, password, httpClient)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	jobs := cron.BuildJobsFromConfig(configPath)

	for _, job := range jobs {
		query := search.New(host, job.Name, job.Option.Query, job.Option.Streamid, job.Frequency, job.Option.Fields, httpClient)

		graylogService := graylog.New(sessionService, query)

		s.Every(job.Frequency).Minutes().Do(job.GetCron(graylogService, report.Report{}))
	}

	s.StartBlocking()
}
