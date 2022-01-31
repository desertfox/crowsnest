package job

import (
	"github.com/desertfox/crowsnest/config"
)

type Job struct {
	Name      string    `yaml:"name"`
	Host      string    `yaml:"host"`
	Frequency int       `yaml:"frequency"`
	Search    Search    `yaml:"search"`
	Condition Condition `yaml:"condition"`
	Output    Output    `yaml:"output"`
	History   *History  `yaml:"-"`
}

type List struct {
	Jobs   []*Job
	Config *config.Config
}

func (j *Job) Func() func() {
	return func() {
		j := j

		rawCSV := j.Search.Run(j.Frequency)

		result := j.Condition.Parse(rawCSV)

		j.History.Push(result)

		j.Output.Send(j.Name, j.Frequency, j.Search, j.Condition, result)
	}
}
