package job

import (
	"github.com/desertfox/crowsnest/config"
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

		j.Search.Run(j.Frequency)

		j.Search.Send(j.Name, j.Frequency)
	}
}
