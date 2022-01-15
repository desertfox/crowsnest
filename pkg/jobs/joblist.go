package jobs

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type JobList []Job

func (jl *JobList) GetConfig(configPath string) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("unable to read file %s", configPath)
	}

	data := make(map[string]*JobList)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		log.Fatalf("unable to load jobs %s", file)
	}

	if _, ok := data["jobs"]; !ok {
		log.Fatalf("missing jobs yaml key %s", file)
	}

	if len(*data["jobs"]) > 0 {
		for i, job := range *data["jobs"] {
			log.Printf("loaded Job %d: %s", i, job.Name)

			jl.Add(job)
		}
	}
}

func (jl JobList) WriteConfig(configPath string) error {
	var list = map[string]JobList{"jobs": jl}
	data, err := yaml.Marshal(&list)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, data, os.FileMode(int(0777)))
	if err != nil {
		return err
	}

	return nil
}

func (jl JobList) checkIfExists(j Job) bool {
	for _, job := range jl {
		if job.Name == j.Name {
			return true
		}
	}
	return false
}

func (jl *JobList) Add(j Job) error {
	if jl.checkIfExists(j) {
		return errors.New("job exists")
	}

	*jl = append(*jl, j)

	return nil
}

func (jl *JobList) Del(name string) JobList {
	jobs := []Job(*jl)

	for i, j := range jobs {
		if j.Name == name {
			jobs[i] = jobs[len(jobs)-1]
			jobs = jobs[:len(jobs)-1]
			return jobs
		}
	}

	return JobList{}
}

func (jl *JobList) HandleEvent(event Event, configPath string) {
	log.Printf("handleEvent, event: %v", event)

	switch event.Action {
	case ReloadJobList:
		jl.GetConfig(configPath)
	case DelJob:
		jl.Del(event.Value)
		jl.WriteConfig(configPath)
	case AddJob:
		jl.Add(event.Job)
		jl.WriteConfig(configPath)
	}
}
