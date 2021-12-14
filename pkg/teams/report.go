package teams

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type reporter interface {
	GetTitle() string
	GetBodyText() string
	GetURLQuery() string
}

type Report struct {
	title string
	text  string
}

func NewReport(reportTitle, reportText string) Report {
	return Report{
		title: reportTitle,
		text:  reportText,
	}
}

func (r Report) sendReport(teamsUrl string) error {
	card := newCard(r.title, r.text)

	json_data, err := json.Marshal(card)
	if err != nil {
		return err
	}

	resp, err := http.Post(teamsUrl, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
