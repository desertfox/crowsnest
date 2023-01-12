package job

import "fmt"

const (
	TeamsBodyTemplate  string = "🔎 Name: %s<br>⌚ Freq: %d<br>🧮 Count: %d<br>📜 Status: %s<br>🔗 Link: [GrayLog](%s)"
	TeamsTitleTemplate string = "<%s>"
)

type Teams struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}
type Output struct {
	Verbose bool  `yaml:"verbose"`
	Teams   Teams `yaml:"teams"`
}

func (o Output) URL() string {
	return o.Teams.Url
}

func (o Output) format(name string, frequency, count int, isAlert, graylogUrl string) string {
	return fmt.Sprintf(TeamsBodyTemplate, name, frequency, count, isAlert, graylogUrl)
}
