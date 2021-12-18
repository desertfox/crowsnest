package main

import (
	"os"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	rp, err := newReqParams(os.Getenv("CROWSNEST_USERNAME"), os.Getenv("CROWSNEST_PASSWORD"), os.Getenv("CROWSNEST_CONFIG"))
	if err != nil {
		bailOut(err)
	}

	c, err := buildConfigFromENV(rp)
	if err != nil {
		bailOut(err)
	}

	s := gocron.NewScheduler(time.UTC)

	for _, j := range c.Jobs {
		s.Every(j.Frequency).Minutes().Do(j.getJob(c))
	}

	s.StartBlocking()
}
