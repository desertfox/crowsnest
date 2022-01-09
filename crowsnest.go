package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/desertfox/crowsnest/pkg/api"
	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/jobs"
	"github.com/desertfox/crowsnest/pkg/teams/report"
	"github.com/fatih/color"
	"github.com/go-co-op/gocron"
)

var (
	httpClient *http.Client      = &http.Client{}
	un         string            = os.Getenv("CROWSNEST_USERNAME")
	pw         string            = os.Getenv("CROWSNEST_PASSWORD")
	configPath string            = os.Getenv("CROWSNEST_CONFIG")
	s          *gocron.Scheduler = gocron.NewScheduler(time.UTC)
	runServer  bool
)

type crowsnest struct {
	jobs jobs.JobList
}

func init() {
	flag.BoolVar(&runServer, "server", false, "--server=true to start server instance")
}

func main() {
	flag.Parse()

	color.Yellow("Crowsnest Startup")

	jobList, err := jobs.BuildFromConfig(configPath)
	if err != nil {
		panic(err.Error())
	}

	color.Yellow("Crowsnest JobRunner Startup")
	crowsnest{jobList}.Run()

	if runServer {
		color.Yellow("Crowsnest Server Startup")
		api.NewServer(&http.ServeMux{}, configPath, jobList).Run()
	}
}

func (cn crowsnest) Run() {
	color.Yellow("Crowsnest ScheduleJobs")
	cn.ScheduleJobs(un, pw)

	color.Green("Crowsnest Daemon...")
	cn.StartAsync()
}

func (cn *crowsnest) ScheduleJobs(un, pw string) {
	for i, j := range cn.jobs {
		sessionService := session.New(
			j.Search.Host,
			un,
			pw,
			httpClient,
		)

		queryService := search.New(
			j.Search.Host,
			j.Search.Query,
			j.Search.Streamid,
			j.Frequency,
			j.Search.Fields,
			httpClient,
		)

		searchService := jobs.SearchService{
			SessionService: sessionService,
			QueryService:   queryService,
		}

		reportService := report.Report{
			Url: j.TeamsURL,
		}

		s.Every(j.Frequency).Minutes().Do(j.GetCron(searchService, reportService))

		color.Green(fmt.Sprintf("Scheduled Job %d: %s", i, j.Name))
	}
}

func (cn crowsnest) StartAsync() {
	s.StartAsync()
}
