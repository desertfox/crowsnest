package job

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
)

var (
	jobHeaders     []string = []string{"Frequency", "Count", "   Status   ", "Graylog Link"}
	jobTemplateStr          = `<table><tr>{{ range .Headers }}<th>{{.}}</th>{{end}}</tr><tr><td>{{ .Frequency }}</td><td>{{ .Count }}</td><td>{{ .Alert }}</td><td>[GrayLog]({{ .Link }})</td></tr></table>`
	jobTemplate    *template.Template
)

func init() {
	jobTemplate, _ = template.New("teams").Parse(jobTemplateStr)
}

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
	//History
	History *History `yaml:"-"`
}

type Teams struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

func (j *Job) GetFunc(g SearchClient, t *goteamsnotify.TeamsClient) func() {
	return func() {
		j := j

		result := j.Search.Run(g, j.Frequency)

		j.History.Add(result)

		if j.Verbose || j.Condition.IsAlert(result) {
			card := messagecard.NewMessageCard()
			card.Title = fmt.Sprintf("Crowsnest: %s", j.Name)

			var b bytes.Buffer
			err := jobTemplate.Execute(&b, struct {
				Headers   []string
				Frequency int
				Count     int
				Alert     string
				Link      string
			}{
				Headers:   jobHeaders,
				Frequency: j.Frequency,
				Count:     result.Count,
				Alert:     j.Condition.IsAlertText(result),
				Link:      j.Search.BuildURL(j.Host, result.From(j.Frequency), result.To()),
			})
			if err != nil {
				log.Println(err)
			}
			card.Text = b.String()

			if err := t.Send(j.Teams.Url, card); err != nil {
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
