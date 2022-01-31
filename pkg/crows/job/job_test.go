package job

import "testing"

func Test_job(t *testing.T) {
	t.Run("Job", func(t *testing.T) {
		j := testJob()

		tests := []struct {
			name string
			job  Job
		}{
			{"Job Run", j},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				tt.job.Func()
			})
		}
	})
}

func testJob() Job {
	return Job{
		Name:      "Test Job",
		Host:      "https://host.com",
		Frequency: 15,
		Search: Search{
			Type:     "relative",
			Streamid: "abcd12345",
			Query:    "error",
			Fields:   []string{"source", "message"},
		},
		Condition: Condition{
			Threshold: 1,
			State:     ">",
		},
		Output: Output{
			Verbose: 1,
			Teams: Teams{
				Url:  "https://mircosoft.com",
				Name: "Room Name",
			},
		},
	}
}
