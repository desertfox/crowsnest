package report

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/teams/cards"
)

type Report struct{}

func (r Report) Send(url, title, text string) error {
	card := cards.NewCard(title, text)

	json_data, err := json.Marshal(card)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
