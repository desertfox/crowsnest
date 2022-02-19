package job

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/desertfox/crowsnest/graylog"
	"github.com/desertfox/crowsnest/teams"
	"gopkg.in/yaml.v2"
)

var JobPath string = os.Getenv("CROWSNEST_CONFIG")

type Lister interface {
	HandleEvent(Event)
	Jobs() []*Job
}

type List struct {
	jobs []*Job
}

func Load() *List {
	var l *List
	file, err := ioutil.ReadFile(JobPath)
	if err != nil {
		log.Printf("error unable to config file: %v, error: %s", JobPath, err)
		return l
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

func (l List) Jobs() []*Job {
	return l.jobs
}

func (l *List) Add(j *Job) error {
	if l.Exists(j) {
		return errors.New("job exists")
	}

	j.Search.Client = graylog.New(j.Host)

	j.Output.Client = teams.Client{}

	j.History = newHistory()

	l.jobs = append(l.jobs, j)

	return nil
}

func (l List) Exists(j *Job) bool {
	for _, job := range l.jobs {
		if job.Name == j.Name {
			return true
		}
	}
	return false
}

func (l List) Save() {
	var jobs = map[string][]*Job{"jobs": l.jobs}
	data, err := yaml.Marshal(&jobs)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(JobPath, data, os.FileMode(int(0777))); err != nil {
		log.Fatal(err)
	}
}

func (l *List) Del(delJob *Job) {
	jobs := []*Job(l.jobs)

	for i, j := range jobs {
		if j.Name == delJob.Name {
			jobs[i] = jobs[len(jobs)-1]
			jobs = jobs[:len(jobs)-1]

			l.jobs = jobs

			break
		}
	}
}

func (l *List) HandleEvent(event Event) {
	switch event.Action {
	case Reload:
		l = Load()
	case Del:
		l.Del(event.Job)
	case Add:
		l.Add(event.Job)
	}
	l.Save()
}
