package job

import "github.com/desertfox/crowsnest/pkg/teams"

type Output struct {
	Verbose  int    `yaml:"verbose"`
	TeamsURL string `yaml:"teamsurl"`
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
