package job

import (
	"fmt"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/config"
	"github.com/desertfox/crowsnest/pkg/crows/job/output"
	"github.com/desertfox/crowsnest/pkg/crows/job/search"
)

var (
	httpClient *http.Client = &http.Client{}
)

type Job struct {
	Name      string        `yaml:"name"`
	Host      string        `yaml:"host"`
	Frequency int           `yaml:"frequency"`
	Search    search.Search `yaml:"search"`
	Output    output.Output `yaml:"output"`
	Config    *config.Config
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
	j.Search.Run(j.Host, j.Config.Username, j.Config.Password, j.Frequency, httpClient)
}

func (j Job) Send() {
	if j.Output.IsVerbose() || j.Search.Condition.IsAlert() {
		j.Output.Send(
			j.Output.URL,
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
