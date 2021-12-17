package main

import (
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	c, err := buildConfigFromENV()
	if err != nil {
		bailOut(err)
	}

	s := gocron.NewScheduler(time.UTC)

	for _, j := range c.Jobs {
		s.Every(j.Frequency).Minutes().Do(j.getJob(c))
	}

	s.StartBlocking()
}
