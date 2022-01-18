package output

import "github.com/desertfox/crowsnest/pkg/teams"

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
	return teams.Report{
		Url: o.TeamsURL,
	}
}
