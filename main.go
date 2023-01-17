package main

import (
	"log"
	"os"

	"github.com/desertfox/crowsnest/api"
	"github.com/desertfox/crowsnest/pkg/crows"
)

const (
	version string = "v3.0"
)

func main() {
	log.Printf("Crowsnest Startup Version %s ", version)

	c := crows.Config{}
	err := c.Load(os.Getenv("CROWSNEST_CONFIG"))
	if err != nil {
		log.Fatal(err)
	}

	nest := c.BuildNest()

	if err := nest.Start(); err != nil {
		panic(err)
	}

	api.New(nest).Start()
}
