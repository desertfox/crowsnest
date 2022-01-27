package job

import (
	"github.com/desertfox/crowsnest/pkg/crows/job/search"
)

func testJob() Job {
	return Job{
		Name:      "Test Job",
		Host:      "https://host.com",
		Frequency: 15,
		Search: search.Search{
			Type:     "relative",
			Streamid: "abcd12345",
			Query:    "error",
			Fields:   []string{"source", "message"},
			Condition: search.Condition{
				Threshold: 1,
				State:     ">",
			},
			Output: search.Output{
				Verbose: 1,
				Teams: search.Teams{
					Url:  "https://mircosoft.com",
					Name: "Room Name",
				},
			},
		},
	}
}
