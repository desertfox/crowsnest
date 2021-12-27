package cards

import "fmt"

type msCard struct {
	Type        string       `json:"type"`
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	ContentType string  `json:"contentType"`
	ContentUrl  string  `json:"contentUrl"`
	Content     content `json:"content"`
}

type content struct {
	Schema  string `json:"$schema"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Body    []body `json:"body"`
}

type body struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Size     string `json:"size"`
	Wrap     bool   `json:"wrap"`
	Markdown bool   `json:"markdown"`
}

func NewCard(reportTitle, reportText string) msCard {
	b := body{
		"TextBlock",
		fmt.Sprintf("%v\n\r%v", reportTitle, reportText),
		"Medium",
		true,
		true,
	}

	a := attachment{
		ContentType: "application/vnd.microsoft.card.adaptive",
		ContentUrl:  "",
		Content: content{
			"http://adaptivecards.io/schemas/adaptive-card.json", "AdaptiveCard", "1.2", []body{b},
		},
	}

	return msCard{
		Type:        "message",
		Attachments: []attachment{a},
	}
}
