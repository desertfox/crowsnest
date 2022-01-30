package job

import (
	"time"

	"github.com/desertfox/crowsnest/config"
)

type Job struct {
	Name      string         `yaml:"name"`
	Host      string         `yaml:"host"`
	Frequency int            `yaml:"frequency"`
	Search    Search         `yaml:"search"`
	Config    *config.Config `yaml:"-"`
	Results   Results        `yaml:"-"`
}

type List struct {
	Jobs   []*Job
	Config *config.Config
}

func (j *Job) Func() func() {
	return func() {
		j := j

		rawCSV := j.Search.Run(j.Frequency)

		count := j.Search.Condition.Parse(rawCSV)

		result := Result{
			Count: count,
			When:  time.Now(),
		}

		j.Results = append(j.Results, result)

		j.Search.Send(j.Name, j.Frequency, result)
	}
}
