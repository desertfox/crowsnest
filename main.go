package main

import (
	"log"

	"github.com/desertfox/crowsnest/pkg/api"
	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/crows"
)

const (
	version string = "v1.6"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	config := config.LoadFromEnv()

	nest := crows.Load(config)
	nest.Run()

	api.New(nest).Run()
}
