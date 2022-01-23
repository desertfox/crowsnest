package job

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type List []*Job

func (jl List) Load(configPath string) List {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("unable to read file %s", configPath)
	}

	data := make(map[string]*List)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		log.Fatalf("unable to load jobs %s", file)
	}

	if _, ok := data["jobs"]; !ok {
		log.Fatalf("missing jobs yaml key %s", file)
	}

	var list List
	if len(*data["jobs"]) > 0 {
		for i, job := range *data["jobs"] {
			log.Printf("loaded Job from config %d: %s", i, job.Name)

			list.Add(job)
		}
	}

	return list
}

func (jl List) Save(configPath string) {
	var list = map[string]List{"jobs": jl}
	data, err := yaml.Marshal(&list)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(configPath, data, os.FileMode(int(0777))); err != nil {
		log.Fatal(err)
	}
}

func (jl *List) Add(j *Job) error {
	if jl.Exists(j) {
		return errors.New("job exists")
	}

	*jl = append(*jl, j)

	return nil
}

func (jl List) Exists(j *Job) bool {
	for _, job := range jl {
		if job.Name == j.Name {
			return true
		}
	}
	return false
}

func (jl *List) Del(delJob *Job) {
	jobs := []*Job(*jl)

	for i, j := range jobs {
		if j.Name == delJob.Name {
			jobs[i] = jobs[len(jobs)-1]
			jobs = jobs[:len(jobs)-1]

			*jl = List(jobs)
			break
		}
	}
}
