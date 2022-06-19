package job

import (
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/desertfox/crowsnest/graylog"
	"github.com/desertfox/crowsnest/teams"
	"gopkg.in/yaml.v3"
)

const (
	Add int = iota
	Del
	Reload
)

var JobPath string = os.Getenv("CROWSNEST_CONFIG")

type Event struct {
	Action int
	Value  string
	Job    *Job
}

type Lister interface {
	HandleEvent(Event)
	Jobs() []*Job
}

type List struct {
	mu   sync.Mutex
	jobs []*Job
}

func NewList() *List {
	var l *List = &List{}

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
		var wg sync.WaitGroup
		for i, job := range data["jobs"] {
			log.Printf("loaded Job from config %d: %s", i, job.Name)
			wg.Add(1)
			go func(job *Job) {
				defer wg.Done()
				l.add(job)
			}(job)
		}
		wg.Wait()
	}

	return l
}

func (l *List) Jobs() []*Job {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.jobs
}

func (l *List) add(j *Job) {
	if l.exists(j) {
		return
	}

	j.Search.Client = graylog.New(j.Host, os.Getenv("CROWSNEST_USERNAME"), os.Getenv("CROWSNEST_PASSWORD"))

	j.Output.Client = teams.Client{}

	j.History = newHistory()

	l.mu.Lock()
	l.jobs = append(l.jobs, j)
	l.mu.Unlock()
}

func (l *List) exists(j *Job) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, job := range l.jobs {
		if job.Name == j.Name {
			return true
		}
	}
	return false
}

func (l *List) save() {
	var jobs = map[string][]*Job{"jobs": l.jobs}
	data, err := yaml.Marshal(&jobs)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(JobPath, data, os.FileMode(int(0777))); err != nil {
		log.Fatal(err)
	}
}

func (l *List) del(delJob *Job) {
	jobs := []*Job(l.jobs)
	for i, j := range jobs {
		if j.Name == delJob.Name {
			jobs[i] = jobs[len(jobs)-1]
			jobs = jobs[:len(jobs)-1]

			l.mu.Lock()
			l.jobs = jobs
			l.mu.Unlock()
			return
		}
	}
}

func (l *List) HandleEvent(event Event) {
	switch event.Action {
	case Reload:
		l = NewList()
	case Del:
		l.del(event.Job)
	case Add:
		l.add(event.Job)
	}
	l.save()
}
