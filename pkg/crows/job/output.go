package job

import (
	"fmt"
)

type Teams struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}
type Output struct {
	Verbose int           `yaml:"verbose"`
	Teams   Teams         `yaml:"teams"`
	Client  OutputService `yaml:"-"`
}
type OutputService interface {
	Send(string, string) error
}

func (o Output) IsVerbose() bool {
	return o.Verbose > 0
}

func (o Output) Send(name string, frequency int, s Search, c Condition, r Result) {
	if o.IsVerbose() || c.IsAlert(r) {
		o.Client.Send(
			o.URL(),
			fmt.Sprintf("🔎 Name: %s<br>"+
				"⌚ Freq: %d<br>"+
				"🧮 Count: %d<br>"+
				"📜 Status: %s<br>"+
				"🔗 Link: [GrayLog](%s)",
				name,
				frequency,
				r.Count,
				c.IsAlertText(r),
				s.BuildURL(r.From(frequency), r.To()),
			),
		)
	}
}

func (o Output) URL() string {
	return o.Teams.Url
}
