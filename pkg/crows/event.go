package crows

import (
	"sync"

	"github.com/desertfox/crowsnest/pkg/crows/job"
)

var (
	loadEventCallbackOnce sync.Once
)

type action int

const (
	Add action = iota
	Del
	Reload
)

type Event struct {
	Action action
	Value  string
	Job    *job.Job
}

func loadEventCallback(n Nest) {
	loadEventCallbackOnce.Do(func() {
		go func() {
			n := n
			event := <-n.eventCallback

			switch event.Action {
			case Reload:
				n.list = &job.List{}
				n.list.Load(n.config)
			case Del:
				n.list.Del(event.Job)
			case Add:
				n.list.Add(event.Job)
			}

			n.list.Save()

			n.scheduler.ClearAndRun(n.list)
		}()
	})
}
