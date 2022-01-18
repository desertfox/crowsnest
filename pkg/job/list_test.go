package job

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_BuildFromConfig(t *testing.T) {
	t.Run("BuildFromConfig", func(t *testing.T) {

		file, err := ioutil.TempFile("", "test_yaml")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())

		fakeyaml := newFakeYaml()
		file.Write(fakeyaml)
		defer file.Close()

		got := List{}

		got = got.Load(file.Name())

		log.Printf("%v", got)

		jobCopy := (got)[0]

		jobCopy.Name = "test"

		gotErr := got.Add(jobCopy)
		if gotErr != nil {
			t.Error(gotErr.Error())
		}

		if len(got) != 2 {
			t.Errorf("wrong number of jobs: %#v", got)
		}

		got.Save(file.Name())

		got.Load(file.Name())

		if len(got) != 2 {
			t.Errorf("wrong number of jobs: %#v", got)
		}

	})
}

func newFakeYaml() []byte {
	return []byte(`jobs:
- name: "DB Errors"
  frequency: 15
  threshold: 20
  teamsurl: ""
  options:
    host: "http://catfacts.com"
    type: "relative"
    streamid: "adfasdfasdf"
    query: "region:production AND DBI AND error"
    fields:
      - "message"
      - "region"
      - "kubernetes_namespace_name"`)
}