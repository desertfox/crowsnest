package cards

type thumbnailCard struct {
	ContentType string           `json:"contentType"`
	Content     thumbnailContent `json:"content"`
}

type thumbnailContent struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	Text     string   `json:"text"`
	Images   []string `json:"images"`
	Buttons  []string `json:"buttons"`
}

func NewThumbnailCard(title, subtitle, text string) thumbnailCard {
	return thumbnailCard{
		ContentType: "application/vnd.microsoft.card.thumbnail",
		Content: thumbnailContent{
			Title:    title,
			Subtitle: subtitle,
			Text:     text,
		},
	}
}
