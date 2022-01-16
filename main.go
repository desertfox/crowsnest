package main

import (
	"log"

	"github.com/desertfox/crowsnest/pkg/api"
	"github.com/desertfox/crowsnest/pkg/jobs"
)

const (
	version string = "v1.2"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	jobs := jobs.Load().Schedule()

	log.Print("Crowsnest API Startup")

	api.New(jobs).Run()
}
