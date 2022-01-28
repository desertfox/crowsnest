package job

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/desertfox/crowsnest/config"
	"github.com/desertfox/crowsnest/graylog"
	"github.com/desertfox/crowsnest/teams"
	"gopkg.in/yaml.v2"
)

var (
	httpClient *http.Client = &http.Client{}
)

type List struct {
	Jobs   []*Job
	Config *config.Config
}

func (l *List) Load() *List {
	file, err := ioutil.ReadFile(l.Config.Path)
	if err != nil {
		log.Printf("error unable to config file: %v, error: %s", l.Config.Path, err)
		return &List{}
	}

	data := make(map[string][]*Job)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		log.Fatalf("unable to load jobs %s", file)
	}

	if _, ok := data["jobs"]; !ok {
		log.Fatalf("missing jobs yaml key %s", file)
	}

	if len(data["jobs"]) > 0 {
		for i, job := range data["jobs"] {
			log.Printf("loaded Job from config %d: %s", i, job.Name)

			l.Add(job)
		}
	}

	return l
}

func (l List) Save() {
	var jobs = map[string][]*Job{"jobs": l.Jobs}
	data, err := yaml.Marshal(&jobs)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(l.Config.Path, data, os.FileMode(int(0777))); err != nil {
		log.Fatal(err)
	}
}

func (l *List) Add(j *Job) error {
	if l.Exists(j) {
		return errors.New("job exists")
	}

	j.Config = l.Config

	j.Search.Client = graylog.New(
		j.Config.Username,
		j.Config.Password,
		j.Host,
		httpClient,
	)

	j.Search.Output.Client = teams.Client{}

	l.Jobs = append(l.Jobs, j)

	return nil
}

func (l List) Exists(j *Job) bool {
	for _, job := range l.Jobs {
		if job.Name == j.Name {
			return true
		}
	}
	return false
}

func (l *List) Del(delJob *Job) {
	jobs := []*Job(l.Jobs)

	for i, j := range jobs {
		if j.Name == delJob.Name {
			jobs[i] = jobs[len(jobs)-1]
			jobs = jobs[:len(jobs)-1]

			l.Jobs = jobs

			break
		}
	}
}

func (l *List) Clear() *List {
	l.Jobs = []*Job{}
	return l
}
