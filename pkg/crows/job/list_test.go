package job

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
)

func Test_List(t *testing.T) {

	t.Run("Save, Load", func(t *testing.T) {
		file, err := os.CreateTemp("", "SaveLoad")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		defer os.Remove(file.Name())

		testJob := testJob("SaveLoad")
		list := List{
			Jobs: []*Job{&testJob},
			File: file.Name(),
		}
		list.Save()

		list = List{
			File: file.Name(),
		}
		list.Load()
		got := len(list.Jobs)

		if got != 1 {
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

/*
func ExampleLoad() {
	job := testJob("example")

	data, err := yaml.Marshal(job)
	if err != nil {
		panic(err)
	}


	fmt.Println(Condition{
		Threshold: 1,
		State:     "<",
	})

}
*/
