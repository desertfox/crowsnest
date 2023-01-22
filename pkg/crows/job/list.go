package job

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type List struct {
	//File path to store jobs
	File string
	//Job array
	Jobs []*Job
	mux  sync.Mutex
}

// Load reads file loaded at List.File and attempts to populate job list
func (l *List) Load() error {
	file, err := os.ReadFile(l.File)
	if err != nil {
		return fmt.Errorf("unable to load list from file, %w", err)
	}

	if err := yaml.Unmarshal(file, &l.Jobs); err != nil {
		return fmt.Errorf("unable to load list from yaml, %w", err)
	}

	for _, j := range l.Jobs {
		j.History = &History{
			results: make([]Result, maxHistory),
		}
	}

	return nil
}
func (l *List) Save() error {
	data, err := yaml.Marshal(l.Jobs)
	if err != nil {
		return err
	}

	if err := os.WriteFile(l.File, data, os.FileMode(int(0644))); err != nil {
		return err
	}

	return nil
}
func (l *List) Count() int {
	return len(l.Jobs)
}
func (l *List) Add(j *Job) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.exists(j) {
		return
	}

	l.Jobs = append(l.Jobs, j)
}
func (l *List) Delete(delJob *Job) {
	l.mux.Lock()
	defer l.mux.Unlock()

	jobs := []*Job(l.Jobs)
	for i, j := range jobs {
		if j.Name == delJob.Name {
			jobs[i] = jobs[len(jobs)-1]
			jobs = jobs[:len(jobs)-1]
			l.Jobs = jobs
			return
		}
	}
}
func (l *List) exists(j *Job) bool {
	for _, job := range l.Jobs {
		if job.Name == j.Name {
			return true
		}
	}
	return false
}
