package main

import (
	"fmt"
	"os"
	"time"
)

func (j job) getJob(q Query, r Report) func() {
	return func() {
		fmt.Println("ExecuteJob " + j.Name)

		count, err := q.Execute()
		if err != nil {
			bailOut(err)
		}

		fmt.Println(time.Now(), count, q.BuildHumanURL())

		var status string
		if count >= j.Threshold {
			status = "ALERT"
		} else {
			status = "OK"
		}

		r.Send(j.Name, fmt.Sprintf("Status: %s\nCount: %d\nLink: [GrayLog Query](%s)\n", status, count, q.BuildHumanURL()))
	}
}

func bailOut(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
