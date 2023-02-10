package job

import (
	"bytes"
	"encoding/csv"
	"errors"
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
	History   *History  `yaml:"-"`
}

type Teams struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

func (j *Job) GetFunc(goclient gograylog.ClientInterface, teamsclient *goteamsnotify.TeamsClient, log *zap.SugaredLogger) func() {
	return func() {
		j := j

		b, err := goclient.Search(j.Search.query)
		if err != nil {
			log.Errorw("unable to complete search", "name", j.Name, "error", err)
			return
		}

		numLines, err := parseCSV(b)
		if err != nil {
			log.Errorw("unable to parse search results", "name", j.Name, "dataPartial", b[64:], "error", err)
			return
		}

		r := &Result{
			Count: numLines,
			When:  time.Now(),
		}

		j.Condition.IsAlert(r)

		j.History.Add(r)

		log.Infow("job run", "name", j.Name, "count", r.Count, "IsAlert", r.Alert, "AlertCount", j.History.alertCount)

		if j.Verbose || r.Alert {
			if err := teamsclient.Send(j.Teams.Url, createTeamsCard(j, r)); err != nil {
				log.Errorw("unable to send results to webhook", "name", j.Name, "error", err)
				return
			}
		}
	}
}

func createTeamsCard(j *Job, r *Result) *messagecard.MessageCard {
	card := messagecard.NewMessageCard()
	card.Title = fmt.Sprintf("Crowsnest: %s", j.Name)
	card.Text = fmt.Sprintf(
		"🔎 Name: %s<br>⌚ Freq: %dm<br>🧮 Count: %d<br>🚨 Alerts: %d<br>📜 Status: %s<br>🔗 Link: [GrayLog](%s)<br><br>[Crowsnest Status Page](%s)",
		j.Name,
		j.Frequency,
		r.Count,
		j.History.alertCount,
		j.Condition.IsAlertText(r),
		j.Search.BuildURL(j.Host, r.From(j.Frequency), r.To()),
		CROWSNEST_STATUS_URL,
	)
	return card
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

func parseCSV(b []byte) (int, error) {
	records, err := csv.NewReader(bytes.NewBuffer(b)).ReadAll()
	if err != nil {
		return 0, err
	}

	if len(records) == 0 {
		return 0, errors.New("response was not text/csv")
	}

	return len(records) - 1, nil
}
