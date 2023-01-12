package crows

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
	"github.com/desertfox/crowsnest/pkg/crows/cron"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/gograylog"
)

var (
	wg sync.WaitGroup
)

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

type Nest struct {
	//List Collection of Jobs
	List *job.List
	//Scheduler controller
	Schedule *cron.Schedule
	//Teams Client for output
	MSTeamsClient *goteamsnotify.TeamsClient
	//Graylog Client for searching
	GrayLogClient *gograylog.Client
}

// Start will check if there are Jobs attached to the list struct value, if not it will attempt to List.Load()
// Schedule will initalize and start schedule control
// Nest method AssignJobs will then be called.
func (n *Nest) Start() error {
	if n.List.Count() == 0 {
		err := n.List.Load()
		if err != nil {
			return err
		}
	}

	n.Schedule.Start()

	n.AssignJobs()

	return nil
}

// AssignJobs will attach all Jobs to the Schedule
// Nest will also attach a Status Job here that is not a Job type but a one off reporter
func (n *Nest) AssignJobs() {
	if n.List.Count() == 0 {
		return
	}

	wg.Add(n.List.Count())
	for _, j := range n.List.Jobs {
		go func(name string, frequency int, startAt time.Time, f func()) {
			defer wg.Done()
			n.Schedule.Add(name, frequency, startAt, f, true)
		}(j.Name, j.Frequency, j.GetOffSetTime(), j.GetFunc(n.GrayLogClient, n.MSTeamsClient))
	}
	wg.Wait()

	n.Schedule.Add("Status Job", 60, time.Now(), n.statusJob(), true)
}

// Jobs Sugar for accessing List.Jobs
func (n *Nest) Jobs() []*job.Job {
	return n.List.Jobs
}

// NextRun will query the Schedule to find the given job name
func (n *Nest) NextRun(name string) time.Time {
	return n.Schedule.NextRun(name)
}

// LastRun will query the Schedule to find the given job name
func (n *Nest) LastRun(name string) time.Time {
	return n.Schedule.LastRun(name)
}

// HandleEvent is the watcher for inbound events from web API to interact
// with running Nest application
func (n *Nest) HandleEvent(event Event) {
	go func(event Event) {
		switch event.Action {
		case Reload:
			n.List = &job.List{
				File: os.Getenv("CROWSNEST_JOBS"),
			}
		case Del:
			n.List.Delete(event.Job)
		case Add:
			n.List.Add(event.Job)
		}
		n.List.Save()

		n.AssignJobs()
	}(event)
}

func (n *Nest) statusJob() func() {
	return func() {
		card := messagecard.NewMessageCard()
		card.Title = "Crowsnest App Status"
		card.Text = fmt.Sprintf("# Jobs Running: %d<br>Uptime: %s", n.List.Count(), "")

		if err := n.MSTeamsClient.Send(os.Getenv("CROWSNEST_TEAMSURL"), card); err != nil {
			log.Printf("unable to send results to webhook %s, %s", "status job", err.Error())
		}

	}
}
