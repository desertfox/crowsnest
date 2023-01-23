package crows

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/atc0005/go-teams-notify/v2/messagecard"
)

var (
	headers           []string = []string{"name", "runs", "avg"}
	statusTemplateStr string   = `<table><tr>{{range .Headers }}<th>{{.}}</th>{{end}}</tr>{{ range .Data }}<tr><td>{{.Name}}</td><td>{{.Avg}}</td></tr>{{end}}</table>`
	statusTemplate    *template.Template
)

type jobRow struct {
	Name string
	Avg  int
}

func init() {
	statusTemplate, _ = template.New("teams").Parse(statusTemplateStr)
}

func (n *Nest) StatusJob() func() {
	return func() {
		var card *messagecard.MessageCard = messagecard.NewMessageCard()
		card.Title = fmt.Sprintf("Crowsnest Status: jobs running %d", n.list.Count())
		card.Text = n.StatusJobOuput()

		if err := n.MSTeamsClient.Send(n.TeamsURL, card); err != nil {
			log.Printf("unable to send results to webhook %s, %s", "status job", err.Error())
		}

	}
}

func (n *Nest) StatusJobOuput() string {
	var (
		b    bytes.Buffer
		rows []jobRow
	)

	for _, j := range n.Jobs() {
		rows = append(rows, jobRow{j.Name, j.History.Avg()})
	}

	if err := statusTemplate.Execute(&b, struct {
		Headers []string
		Data    []jobRow
	}{
		Headers: headers,
		Data:    rows,
	}); err != nil {
		log.Println(err)
		return ""
	}

	return b.String()
}
