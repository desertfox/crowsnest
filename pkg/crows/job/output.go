package job

import "fmt"

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
			fmt.Sprintf("🔎 Name  : %s\n\r"+
				"⌚ Freq  : %d\n\r"+
				"📜 Status: %s\n\r"+
				"🧮 Count : %d\n\r"+
				"🔗 Link  : [GrayLog](%s)",
				name,
				frequency,
				c.IsAlertText(r),
				r.Count,
				s.BuildURL(),
			),
		)
	}
}

func (o Output) URL() string {
	return o.Teams.Url
}
