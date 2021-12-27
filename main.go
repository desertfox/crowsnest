package main

import (
	"net/http"
	"os"

	crowsnest "github.com/desertfox/crowsnest/pkg"
)

var (
	httpClient *http.Client = &http.Client{}
	host       string       = os.Getenv("CROWSNEST_HOST")
	username   string       = os.Getenv("CROWSNEST_USERNAME")
	password   string       = os.Getenv("CROWSNEST_PASSWORD")
	configPath string       = os.Getenv("CROWSNEST_CONFIG")
)

func main() {
	cn := crowsnest.New(
		host,
		username,
		password,
		configPath,
		httpClient,
	)

	cn.ScheduleJobs()

	cn.StartBlocking()
}
