package main

import (
	"log"

	"github.com/desertfox/crowsnest/pkg/api"
	"github.com/desertfox/crowsnest/pkg/scheduler"
)

const (
	version string = "v1.4"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)
	api.New(scheduler.Load().Schedule()).Run()
}
