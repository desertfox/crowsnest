package jobs

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type JobList []Job

func BuildFromConfig(configPath string) []Job {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err.Error())
	}

	data := make(map[string][]Job)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		panic(err.Error())
	}

	return data["jobs"]
}

func AddToConfig(configPath string, job Job) error {
	var jl JobList = BuildFromConfig(configPath)

	if !jl.checkIfExists(job) {
		return errors.New("job exists")
	}

	jl = append(jl, job)

	if err := jl.writeConfig(configPath); err != nil {
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

func (jl JobList) writeConfig(configPath string) error {
	data, err := yaml.Marshal(jl)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, data, os.FileMode(int(0777)))
	if err != nil {
		return err
	}

	return nil
}
