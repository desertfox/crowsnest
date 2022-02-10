package job

import (
	"testing"
	"time"
)

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

func Test_OffSet(t *testing.T) {
	t.Run("Job", func(t *testing.T) {
		j := testJob()

		now := time.Now()

		tests := []struct {
			name       string
			timeString string
			hour       int
			min        int
		}{
			{"No OffSet", "", now.Hour(), now.Minute()},
			{"AM", "01:33", 01, 33},
			{"PM", "13:30", 13, 30},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				j.Offset = tt.timeString
				offset := j.GetOffSetTime()

				if offset.Hour() != tt.hour {
					t.Fatalf("OffSet hour does not match expected. offset:%v, got:%v, want:%v", offset, offset.Hour(), tt.hour)
				}

				if offset.Minute() != tt.min {
					t.Fatalf("OffSet hour does not match expected. offset:%v, got:%v, want:%v", offset, offset.Minute(), tt.min)
				}
			})
		}
	})
}

func testJob() Job {
	return Job{
		Name:      "Test Job",
		Host:      "https://host.com",
		Frequency: 15,
		Offset:    "13:00",
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
