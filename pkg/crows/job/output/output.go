package output

import "github.com/desertfox/crowsnest/pkg/crows/job/output/teams"

type Output struct {
	Verbose int    `yaml:"verbose"`
	URL     string `yaml:"teamsurl"`
}
type OutputService interface {
	Send(string, string) error
}

func (o Output) IsVerbose() bool {
	return o.Verbose > 0
}

func (o Output) Service() OutputService {
	return teams.Report{}
}

func (o Output) Send(url string, text string) {
	o.Service().Send(url, text)
}
