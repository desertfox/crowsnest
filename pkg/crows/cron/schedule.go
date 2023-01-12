package cron

import (
	"errors"
	"time"

	"github.com/go-co-op/gocron"
)

type Schedule struct {
	Gocron *gocron.Scheduler
}

func (s *Schedule) Start() {
	s.Gocron.StartAsync()
}

func (s *Schedule) Add(name string, frequency int, startAt time.Time, do func(), replaceExisting bool) {
	existingJob, err := s.get(name)
	if err == nil && existingJob.IsRunning() {
		if !replaceExisting {
			return
		}
		s.Gocron.Remove(existingJob)
	}
	s.Gocron.Every(frequency).Minutes().StartAt(startAt).Tag(name).Do(do)
}
func (s *Schedule) NextRun(name string) time.Time {
	j, err := s.get(name)
	if err != nil {
		return time.Now()
	}
	return j.NextRun()
}
func (s *Schedule) LastRun(name string) time.Time {
	j, err := s.get(name)
	if err != nil {
		return time.Now()
	}
	return j.LastRun()
}
func (s *Schedule) get(tag string) (*gocron.Job, error) {
	for _, cj := range s.Gocron.Jobs() {
		for _, t := range cj.Tags() {
			if tag == t {
				return cj, nil
			}
		}
	}
	return &gocron.Job{}, errors.New("no job found for tag:" + tag)
}
