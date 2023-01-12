package job

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
)

type Job struct {
	//Name of the job
	Name string `yaml:"name"`
	//Host is the graylog endpoint
	Host string `yaml:"host"`
	//Frequency is the occurence of the job execution
	Frequency int `yaml:"frequency"`
	//Offset if a job is to no begin on startup but at a defered time
	Offset string `yaml:"offset"`
	//Search
	Search Search `yaml:"search"`
	//Condition
	Condition Condition `yaml:"condition"`
	//Output
	Output Output `yaml:"output"`
	//History
	History *History `yaml:"-"`
}

func (j *Job) GetFunc(g SearchClient, t *goteamsnotify.TeamsClient) func() {
	return func() {
		j := j

		result := j.Search.Run(g, j.Frequency)

		j.History.Add(result)

		if j.Output.Verbose || j.Condition.IsAlert(result) {
			card := messagecard.NewMessageCard()
			card.Title = fmt.Sprintf(TeamsTitleTemplate, "Crowsnest")
			card.Text = j.Output.format(j.Name, j.Frequency, result.Count, j.Condition.IsAlertText(result), j.Search.BuildURL(j.Host, result.From(j.Frequency), result.To()))

			if err := t.Send(j.Output.URL(), card); err != nil {
				log.Panicf("unable to send results to webhook %s, %s", j.Name, err.Error())
			}

		}
	}
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
