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
	statusTemplateStr string   = `<table><tr>{{range .Headers }}<th>{{.}}</th>{{end}}</tr>{{ range .Data }}<tr><td>{{.Name}}</td><td>{{.Runs}}</td><td>{{.Avg}}</td></tr>{{end}}</table>`
	statusTemplate    *template.Template
)

type jobRow struct {
	Name string
	Runs int
	Avg  int
}

func init() {
	statusTemplate, _ = template.New("teams").Parse(statusTemplateStr)
}

func (n *Nest) statusJob() func() {
	return func() {
		var (
			n    *Nest = n
			b    bytes.Buffer
			rows []jobRow = make([]jobRow, n.list.Count())
		)

		for _, j := range n.Jobs() {
			rows = append(rows, jobRow{j.Name, len(j.History.Results()), j.History.Avg()})
		}

		if err := statusTemplate.Execute(&b, struct {
			Headers []string
			Data    []jobRow
		}{
			Headers: headers,
			Data:    rows,
		}); err != nil {
			log.Println(err)
			return
		}

		card := messagecard.NewMessageCard()
		card.Title = fmt.Sprintf("Status Job, Running: %d", n.list.Count())
		card.Text = b.String()

		if err := n.MSTeamsClient.Send(n.TeamsURL, card); err != nil {
			log.Printf("unable to send results to webhook %s, %s", "status job", err.Error())
		}

	}
}
