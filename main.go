package main

import (
	"log"

	"github.com/desertfox/crowsnest/pkg/api"
	"github.com/desertfox/crowsnest/pkg/scheduler"
)

const (
	version string = "v1.3"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	s := scheduler.Load().Schedule()

	log.Print("Crowsnest API Startup")

	api.New(s).Run()
}
