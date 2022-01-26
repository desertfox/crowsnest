package job

import (
	"fmt"

	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/crows/job/search"
)

type Job struct {
	Name      string         `yaml:"name"`
	Host      string         `yaml:"host"`
	Frequency int            `yaml:"frequency"`
	Search    search.Search  `yaml:"search"`
	Config    *config.Config `yaml:"-"`
}

// S(un,pw) -> C(S(un,pw)) -> O(url)

func (j Job) Func() func() {
	return func() {
		j := j

		j.Run()

		j.Send()
	}
}

func (j Job) Run() {
	j.Search.Run(j.Frequency)
}

func (j Job) Send() {
	if j.Search.Output.IsVerbose() || j.Search.Condition.IsAlert() {
		j.Search.Output.Send(
			j.Search.Output.URL,
			j.buildText(j.Search.BuildURL()),
		)
	}
}

func (j Job) buildText(url string) string {
	return fmt.Sprintf("🔎 Name  : %s\n\r"+
		"⌚ Freq  : %d\n\r"+
		"📜 Status: %s\n\r"+
		"🧮 Count : %d\n\r"+
		"🔗 Link  : [GrayLog](%s)",
		j.Name,
		j.Frequency,
		j.Search.Condition.IsAlertText(),
		j.Search.Condition.Count,
		url,
	)
}
