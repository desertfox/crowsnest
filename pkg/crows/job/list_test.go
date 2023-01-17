package job

import (
	"fmt"
	"log"
	"os"
	"sync"
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

		testJob := testJob("LoadSave")
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

	t.Run("Load.Add", func(t *testing.T) {
		file, err := os.CreateTemp("", "LoadAdd")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		defer os.Remove(file.Name())

		got := List{}
		var wg sync.WaitGroup
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func(d int) {
				job := testJob(fmt.Sprintf("%d", d))
				got.Add(&job)
				wg.Done()
			}(i)
		}
		wg.Wait()
		want := 10

		if len(got.Jobs) != want {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got.Jobs, want)
		}
	})

	t.Run("Load.Del", func(t *testing.T) {
		got := List{}
		var wg sync.WaitGroup
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func(d int) {
				job := testJob(fmt.Sprintf("%d", d))
				got.Add(&job)
				wg.Done()
			}(i)
		}
		wg.Wait()

		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func(d int) {
				job := testJob(fmt.Sprintf("%d", d))
				got.Delete(&job)
				wg.Done()
			}(i)
		}
		wg.Wait()

		if len(got.Jobs) != 0 {
			t.Errorf("wrong number of jobs, got:%v, want:%v", got.Jobs, "0")
		}
	})

	t.Run("Load.Exists", func(t *testing.T) {
		jobExample := testJob("Exists")
		got := List{
			Jobs: []*Job{&jobExample},
		}

		if !got.exists(&jobExample) {
			t.Errorf("duplicate job returned true, got:%v, want:%v", got.exists(&jobExample), false)
		}
	})
}
