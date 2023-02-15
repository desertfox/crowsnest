package crows

import (
	"sync"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/desertfox/crowsnest/pkg/crows/cron"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/gograylog"
	"go.uber.org/zap"
)

var (
	wg sync.WaitGroup
)

type Nest struct {
	*Config
	list     *job.List
	schedule *cron.Schedule
	//Teams Client for output
	MSTeamsClient *goteamsnotify.TeamsClient
	//Graylog Client for searching
	GrayLogClient gograylog.ClientInterface
	//Event channel for sending signals to an instance
	eventChan chan Event
	log       *zap.SugaredLogger
}

// Start will check if there are Jobs attached to the list struct value, if not it will attempt to List.Load()
// Schedule will initalize and start schedule control
// Nest method AssignJobs will then be called.
func (n *Nest) Start() error {
	if n.list.Count() == 0 {
		n.log.Info("Job list empty")
		err := n.list.Load()
		if err != nil {
			return err
		}
	}

	n.log.Info("schedule started")
	n.schedule.Start()

	n.log.Info("assign jobs")
	n.AssignJobs()

	go handleEvent(n)

	return nil
}

// AssignJobs will attach all Jobs to the Schedule
// Nest will also attach a Status Job here that is not a Job type but a one off reporter
func (n *Nest) AssignJobs() {
	if n.list.Count() == 0 {
		return
	}

	n.log.Infow("assigning jobs", "job count", n.list.Count())
	wg.Add(n.list.Count())
	for _, j := range n.list.Jobs {
		go func(name string, frequency int, startAt time.Time, f func()) {
			defer wg.Done()

			n.log.Infow("adding job", "name", name, "frequency", frequency, "startAt", startAt)
			n.schedule.Add(name, frequency, startAt, f, true)

		}(j.Name, j.Frequency, j.GetOffSetTime(), j.GetFunc(n.GrayLogClient, n.MSTeamsClient, n.log))
	}
	wg.Wait()

}

// Jobs Sugar for accessing List.Jobs
func (n *Nest) Jobs() []*job.Job {
	return n.list.Jobs
}

// NextRun will query the Schedule to find the given job name
func (n *Nest) NextRun(name string) time.Time {
	return n.schedule.NextRun(name)
}

// LastRun will query the Schedule to find the given job name
func (n *Nest) LastRun(name string) time.Time {
	return n.schedule.LastRun(name)
}

func (n *Nest) Send(e Event) {
	n.eventChan <- e
	close(n.eventChan)
}
