package job

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_Load(t *testing.T) {

	t.Run("Load.Save", func(t *testing.T) {
		file, err := ioutil.TempFile("", "LoadSave")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		defer os.Remove(file.Name())

		testJob := testJob()
		JobPath = file.Name()
		list := List{
			jobs: []*Job{&testJob},
		}
		list.Save()

		got, err := ioutil.ReadFile(file.Name())
		if err != nil {
			t.Error(err)
		}

		if len(got) == 1 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, "<1")
		}
	})

	t.Run("Load", func(t *testing.T) {
		file, err := ioutil.TempFile("", "Load")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		defer os.Remove(file.Name())

		testJob := testJob()
		JobPath = file.Name()
		list := List{
			jobs: []*Job{&testJob},
		}
		list.Save()

		got := NewList()
		want := 1

		if len(got.Jobs()) != want {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, want)
		}

		got.Save()
		got = NewList()

		if len(got.Jobs()) != want {
			t.Errorf("wrong number of jobs: %#v", got)
		}

	})

	t.Run("Load.Add", func(t *testing.T) {
		file, err := ioutil.TempFile("", "LoadAdd")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		defer os.Remove(file.Name())

		got := List{}
		job := testJob()
		got.Add(&job)
		want := 1

		if len(got.Jobs()) != want {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, want)
		}
	})

	t.Run("Load.Del", func(t *testing.T) {
		jobExample := testJob()
		got := List{
			jobs: []*Job{&jobExample},
		}

		if len(got.Jobs()) != 1 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, "0")
		}

		got.Del(&jobExample)

		if len(got.Jobs()) != 0 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, "0")
		}
	})

	t.Run("Load.Exists", func(t *testing.T) {
		jobExample := testJob()
		got := List{
			jobs: []*Job{&jobExample},
		}

		if !got.Exists(&jobExample) {
			t.Errorf("duplicate job returned true, got:%v, want:%v", got.Exists(&jobExample), false)
		}
	})
}
