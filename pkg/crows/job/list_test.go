package job

import (
	"log"
	"os"
	"testing"
)

func Test_List(t *testing.T) {

	t.Run("Load.Save", func(t *testing.T) {
		file, err := os.CreateTemp("", "LoadSave")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		defer os.Remove(file.Name())

		testJob := testJob()
		list := List{
			Jobs: []*Job{&testJob},
			File: file.Name(),
		}
		list.Save()

		got, err := os.ReadFile(file.Name())
		if err != nil {
			t.Error(err)
		}

		if len(got) == 1 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, "<1")
		}
	})

	t.Run("Load", func(t *testing.T) {
		file, err := os.CreateTemp("", "Load")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		defer os.Remove(file.Name())

		testJob := testJob()
		list := List{
			Jobs: []*Job{&testJob},
			File: file.Name(),
		}
		list.Save()

		got := List{
			File: file.Name(),
		}
		got.Load()
		want := 1

		if len(got.Jobs) != want {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got, want)
		}

	})

	t.Run("Load.Add", func(t *testing.T) {
		file, err := os.CreateTemp("", "LoadAdd")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		defer os.Remove(file.Name())

		got := List{}
		job := testJob()
		got.Add(&job)
		want := 1

		if len(got.Jobs) != want {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got.Jobs, want)
		}
	})

	t.Run("Load.Del", func(t *testing.T) {
		jobExample := testJob()
		got := List{
			Jobs: []*Job{&jobExample},
		}

		if len(got.Jobs) != 1 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got.Jobs, "0")
		}

		got.Delete(&jobExample)

		if len(got.Jobs) != 0 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got.Jobs, "0")
		}
	})

	t.Run("Load.Exists", func(t *testing.T) {
		jobExample := testJob()
		got := List{
			Jobs: []*Job{&jobExample},
		}

		if !got.exists(&jobExample) {
			t.Errorf("duplicate job returned true, got:%v, want:%v", got.exists(&jobExample), false)
		}
	})
}
