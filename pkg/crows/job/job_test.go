package job

import (
	"sync"
	"testing"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/desertfox/gograylog"
	"go.uber.org/zap"
)

func Test_job(t *testing.T) {
	t.Run("Job", func(t *testing.T) {
		j := testJob("Job")

		tests := []struct {
			name string
			job  Job
		}{
			{"Job Run", j},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				tt.job.GetFunc(&gograylog.Client{}, goteamsnotify.NewTeamsClient(), &zap.SugaredLogger{}, &sync.Mutex{})
			})
		}
	})
}

func Test_OffSet(t *testing.T) {
	t.Run("Job", func(t *testing.T) {
		j := testJob("OffSet")

		now := time.Now()

		tests := []struct {
			name       string
			timeString string
			hour       int
			min        int
		}{
			{"No OffSet", "", now.Hour(), now.Minute() + 1}, //smell
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
					t.Fatalf("OffSet minute does not match expected. offset:%v, got:%v, want:%v", offset, offset.Minute(), tt.min)
				}
			})
		}
	})
}

func testJob(postFixTag string) Job {
	return Job{
		Name:      "Test Job " + postFixTag,
		Frequency: 15,
		Offset:    "13:00",
		Verbose:   true,
		Teams: Teams{
			Url:  "https://mircosoft.com",
			Name: "Room Name",
		},
		Search: Search{
			Streamid: "abcd12345",
			Query:    "error",
			Fields:   []string{"source", "message"},
		},
		Condition: Condition{
			Threshold: 1,
			Operator:  ">",
		},
	}
}
