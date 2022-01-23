package teams

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Report struct{}

func (r Report) Send(url, text string) error {
	card := newCard(text)

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
