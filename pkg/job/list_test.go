package job

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_Load(t *testing.T) {
	t.Run("Load", func(t *testing.T) {
		jobFile := testLoad()
		defer os.Remove(jobFile.Name())

		emptyList := List{}
		got := emptyList.Load(jobFile.Name())
		want := 1

		if len(got) != want {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, want)
		}

		got.Save(jobFile.Name())
		got.Load(jobFile.Name())

		if len(got) != want {
			t.Errorf("wrong number of jobs: %#v", got)
		}

	})

	t.Run("Load.Save", func(t *testing.T) {
		file, err := ioutil.TempFile("", "LoadSave")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		defer os.Remove(file.Name())

		list := List{}
		list.Save(file.Name())

		got, err := ioutil.ReadFile(file.Name())
		if err != nil {
			t.Error(err)
		}

		if len(got) < 1 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, "<1")
		}
	})

	t.Run("Load.Add", func(t *testing.T) {
		jobExample := testJob()

		got := List{}
		got.Add(&jobExample)
		want := 1

		if len(got) != want {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, want)
		}
	})

	t.Run("Load.Del", func(t *testing.T) {
		jobExample := testJob()
		got := List([]*Job{&jobExample})

		if len(got) != 1 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, "0")
		}

		got.Del(&jobExample)

		if len(got) != 0 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, "0")
		}
	})

	t.Run("Load.Exists", func(t *testing.T) {
		jobExample := testJob()
		got := List([]*Job{&jobExample})

		if !got.Exists(&jobExample) {
			t.Errorf("duplicate job returned true, got:%v, want:%v", got.Exists(&jobExample), false)
		}
	})
}

func testLoad() *os.File {
	file, err := ioutil.TempFile("", "test_yaml")
	if err != nil {
		log.Fatal(err)
	}

	fakeyaml := testJobYaml()
	file.Write(fakeyaml)
	defer file.Close()

	return file
}

func testJobYaml() []byte {
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
