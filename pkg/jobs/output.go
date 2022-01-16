package jobs

import "github.com/desertfox/crowsnest/pkg/teams/report"

type Output struct {
	Verbose  int    `yaml:"verbose"`
	TeamsURL string `yaml:"teamsurl"`
}
type ReportService interface {
	Send(string) error
}

func (o Output) isVerbose() bool {
	return o.Verbose > 0
}

func (o Output) reportService() ReportService {
	return report.Report{
		Url: o.TeamsURL,
	}
}
