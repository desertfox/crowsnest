package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/desertfox/crowsnest/api"
	"github.com/desertfox/crowsnest/pkg/crows"
	"github.com/desertfox/crowsnest/pkg/crows/cron"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/gograylog"
	"github.com/go-co-op/gocron"
)

const (
	version string = "v3.0"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	graylog := &gograylog.Client{
		Host: os.Getenv("CROWSNEST_HOST"),
		HttpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
	err := graylog.Login(os.Getenv("CROWSNEST_USERNAME"), os.Getenv("CROWSNEST_PASSWORD"))
	if err != nil {
		log.Fatalf(err.Error())
	}

	nest := &crows.Nest{
		List: &job.List{
			File: os.Getenv("CROWSNEST_JOBS"),
		},
		Schedule: &cron.Schedule{
			Gocron: gocron.NewScheduler(time.UTC),
		},
		MSTeamsClient: goteamsnotify.NewTeamsClient(),
		GrayLogClient: graylog,
	}

	if err := nest.Start(); err != nil {
		panic(err)
	}

	api.New(nest).Start()
}
