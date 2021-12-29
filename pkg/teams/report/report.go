package report

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/teams/cards"
)

type Report struct {
	Url string
}

func (r Report) Send(title, subtitle, text string) error {
	//card := cards.NewCard(title, text)
	card := cards.NewThumbnailCard(title, subtitle, text)

	json_data, err := json.Marshal(card)
	if err != nil {
		return err
	}

	resp, err := http.Post(r.Url, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
