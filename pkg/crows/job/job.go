package job

import (
	"strconv"
	"strings"
	"time"

	"github.com/desertfox/crowsnest/config"
)

type Job struct {
	Name      string    `yaml:"name"`
	Host      string    `yaml:"host"`
	Frequency int       `yaml:"frequency"`
	Offset    string    `yaml:"offset"`
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

func (j Job) HasOffset() bool {
	return j.Offset != ""
}

func (j Job) GetOffSetTime() time.Time {
	if j.Offset == "" {
		return time.Now()
	}
	offSet := strings.Split(j.Offset, ":")
	hour, _ := strconv.Atoi(offSet[0])
	min, _ := strconv.Atoi(offSet[1])

	today := time.Now()
	return time.Date(today.Year(), today.Month(), today.Day(), hour, min, 0, 0, time.UTC)
}
