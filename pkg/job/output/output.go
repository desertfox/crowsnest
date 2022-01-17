package output

import "github.com/desertfox/crowsnest/pkg/teams/report"

type Output struct {
	Verbose  int    `yaml:"verbose"`
	TeamsURL string `yaml:"teamsurl"`
}
type Service interface {
	Send(string) error
}

func (o Output) IsVerbose() bool {
	return o.Verbose > 0
}

func (o Output) Service() Service {
	return report.Report{
		Url: o.TeamsURL,
	}
}
