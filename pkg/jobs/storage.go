package jobs

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func BuildFromConfig(configPath string) (*JobList, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return &JobList{}, err
	}

	data := make(map[string]*JobList)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return &JobList{}, err
	}

	for i, job := range *data["jobs"] {
		log.Printf("Loaded Job %d: %s", i, job.Name)
	}

	return data["jobs"], nil
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
