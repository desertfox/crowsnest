package job

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
	"github.com/desertfox/gograylog"
	"go.uber.org/zap"
)

var CROWSNEST_STATUS_URL string = os.Getenv("CROWSNEST_STATUS_URL")

type Job struct {
	//Name of the job
	Name string `yaml:"name"`
	//Host is the graylog endpoint
	Host string `yaml:"host"`
	//Frequency is the occurence of the job execution
	Frequency int `yaml:"frequency"`
	//Verbose If true will send message to teams room regardless of condition eval
	Verbose bool `yaml:"verbose"`
	//Teams
	Teams Teams `yaml:"teams"`
	//Offset if a job is to no begin on startup but at a defered time
	Offset string `yaml:"offset"`
	//Search
	Search Search `yaml:"search"`
	//Condition
	Condition Condition `yaml:"condition"`
}

type Teams struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

func (j *Job) GetFunc(g gograylog.ClientInterface, t *goteamsnotify.TeamsClient, log *zap.SugaredLogger) func() {
	return func() {
		j := j

		r, err := j.Search.Run(g, j.Frequency)
		if err != nil {
			log.Errorw("unable to compelte search", "name", j.Name, "error", err)
		}

		log.Infow("job run", "name", j.Name, "count", r.Count, "IsAlert", j.Condition.IsAlert(r))

		if j.Verbose || j.Condition.IsAlert(r) {
			if err := t.Send(j.Teams.Url, createTeamsCard(j, r)); err != nil {
				log.Errorw("unable to send results to webhook", "name", j.Name, "error", err)
			}
		}
	}
}

func createTeamsCard(j *Job, r Result) *messagecard.MessageCard {
	card := messagecard.NewMessageCard()
	card.Title = fmt.Sprintf("[Crowsnest](%s): %s", CROWSNEST_STATUS_URL, j.Name)
	card.Text = fmt.Sprintf(
		"🔎 Name: %s<br>⌚ Freq: %d<br>🧮 Count: %d<br>🚨 Alerts: %d<br>📜 Status: %s<br>🔗 Link: [GrayLog](%s)",
		j.Name,
		j.Frequency,
		r.Count,
		j.Alerting(),
		j.Condition.IsAlertText(r),

		j.Search.BuildURL(j.Host, r.From(j.Frequency), r.To()),
	)
	return card
}

// Alerting returns the count of Condition.IsAlert of History Results until it finds a non alerting result
func (j Job) Alerting() int {
	var count int = 0
	for _, r := range j.Search.History.Results() {
		if j.Condition.IsAlert(r) {
			count++
			continue
		}
		break
	}
	return count
}

func (j Job) GetOffSetTime() time.Time {
	today := time.Now()
	if j.Offset == "" {
		return today.Add(1 * time.Minute)
	}

	offSet := strings.Split(j.Offset, ":")
	hour, _ := strconv.Atoi(offSet[0])
	min, _ := strconv.Atoi(offSet[1])

	return time.Date(today.Year(), today.Month(), today.Day(), hour, min, 0, 0, time.UTC)
}
