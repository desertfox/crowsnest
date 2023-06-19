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

	for i := 0; i < len(l.Jobs); i++ {
		l.Jobs[i].History = newHistory()
		l.Jobs[i].Search.buildQuery(l.Jobs[i].Frequency)
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

func (l *List) Add(j *Job) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.exists(j) {
		return
	}

	j.History = newHistory()
	j.Search.buildQuery(j.Frequency)

	l.Jobs = append(l.Jobs, j)
}

func (l *List) Delete(delJob *Job) {
	l.mux.Lock()
	defer l.mux.Unlock()

	for i := 0; i < len(l.Jobs); i++ {
		if l.Jobs[i].Name == delJob.Name {
			l.Jobs[i] = l.Jobs[len(l.Jobs)-1]
			l.Jobs = l.Jobs[:len(l.Jobs)-1]
			return
		}
	}
}

func (l *List) exists(j *Job) bool {
	for i := 0; i < len(l.Jobs); i++ {
		if l.Jobs[i].Name == j.Name {
			return true
		}
	}

	return false
}
