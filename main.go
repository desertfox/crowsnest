package main

import (
	"net/http"
	"os"

	crowsnest "github.com/desertfox/crowsnest/pkg"
	"github.com/fatih/color"
)

var (
	httpClient *http.Client = &http.Client{}
	un         string       = os.Getenv("CROWSNEST_USERNAME")
	pw         string       = os.Getenv("CROWSNEST_PASSWORD")
	configPath string       = os.Getenv("CROWSNEST_CONFIG")
)

func main() {
	color.Yellow("Crowsnest Startup")

	cn := crowsnest.New(
		configPath,
		httpClient,
	)

	color.Yellow("Crowsnest ScheduleJobs")

	cn.ScheduleJobs(un, pw)

	color.Green("Crowsnest Daemon...")

	cn.StartBlocking()
}
