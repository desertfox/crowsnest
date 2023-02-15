package crows

import "github.com/desertfox/crowsnest/pkg/crows/job"

const (
	Add int = iota
	Del
	Reload
)

type Event struct {
	Action int
	Value  string
	Job    *job.Job
}

func handleEvent(n *Nest) {
	for e := range n.eventChan {
		switch e.Action {
		case Reload:
			n.list = n.BuildList()
		case Del:
			n.list.Delete(e.Job)
		case Add:
			n.list.Add(e.Job)
		}
		n.list.Save()

		n.AssignJobs()
	}
	close(n.eventChan)
}
