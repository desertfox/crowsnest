package main

import (
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	s := gocron.NewScheduler(time.UTC)

	c := buildConfigFromENV()

	for _, j := range c.Jobs {
		s.Every(j.getFrequency()).Minutes().Do(j.getFunc(c))
	}

	s.StartBlocking()
}
