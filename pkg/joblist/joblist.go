package joblist

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/desertfox/crowsnest/pkg/job"
	"gopkg.in/yaml.v2"
)

type JobList []job.Job

func (jl *JobList) Load(configPath string) JobList {
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

	var jobList JobList
	if len(*data["jobs"]) > 0 {
		for i, job := range *data["jobs"] {
			log.Printf("loaded Job from config %d: %s", i, job.Name)

			jobList.Add(job)
		}
	}

	return jobList
}

func (jl JobList) Save(configPath string) {
	var list = map[string]JobList{"jobs": jl}
	data, err := yaml.Marshal(&list)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(configPath, data, os.FileMode(int(0777))); err != nil {
		log.Fatal(err)
	}
}

func (jl *JobList) Add(j job.Job) error {
	if jl.Exists(j) {
		return errors.New("job exists")
	}

	*jl = append(*jl, j)

	return nil
}

func (jl JobList) Exists(j job.Job) bool {
	for _, job := range jl {
		if job.Name == j.Name {
			return true
		}
	}
	return false
}

func (jl *JobList) Del(name string) {
	jobs := []job.Job(*jl)

	for i, j := range jobs {
		if j.Name == name {
			jobs[i] = jobs[len(jobs)-1]
			jobs = jobs[:len(jobs)-1]

			*jl = JobList(jobs)
		}
	}
}
