package main

import (
	"net/http"
	"os"

	crowsnest "github.com/desertfox/crowsnest/pkg"
)

var (
	httpClient *http.Client = &http.Client{}
	configPath string       = os.Getenv("CROWSNEST_CONFIG")
)

func main() {
	cn := crowsnest.New(
		configPath,
		httpClient,
	)

	cn.ScheduleJobs()

	cn.StartBlocking()
}
