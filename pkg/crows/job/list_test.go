package job

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/desertfox/crowsnest/pkg/config"
)

func Test_Load(t *testing.T) {
	t.Run("Load", func(t *testing.T) {
		jobFile := testLoad()
		defer os.Remove(jobFile.Name())

		emptyList := List{}
		got := emptyList.Load(&config.Config{
			Path: jobFile.Name(),
		})
		want := 1

		if len(got.Jobs) != want {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, want)
		}

		got.Save()
		got.Load(&config.Config{
			Path: jobFile.Name(),
		})

		if len(got.Jobs) != want {
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

		list := List{
			Config: &config.Config{
				Path: file.Name(),
			},
		}
		list.Save()

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

		if len(got.Jobs) != want {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, want)
		}
	})

	t.Run("Load.Del", func(t *testing.T) {
		jobExample := testJob()
		got := List{
			Jobs: []*Job{&jobExample},
		}

		if len(got.Jobs) != 1 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, "0")
		}

		got.Del(&jobExample)

		if len(got.Jobs) != 0 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, "0")
		}
	})

	t.Run("Load.Exists", func(t *testing.T) {
		jobExample := testJob()
		got := List{
			Jobs: []*Job{&jobExample},
		}

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
