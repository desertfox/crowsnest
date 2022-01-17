package main

import (
	"log"

	"github.com/desertfox/crowsnest/pkg/api"
	"github.com/desertfox/crowsnest/pkg/job"
)

const (
	version string = "v1.4"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)
	api.New(job.Load().Schedule()).Run()
}
