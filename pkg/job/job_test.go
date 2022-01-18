package job

import (
	"github.com/desertfox/crowsnest/pkg/job/condition"
	"github.com/desertfox/crowsnest/pkg/job/output"
	"github.com/desertfox/crowsnest/pkg/job/search"
)

func testJob() Job {
	return Job{
		Name: "Test Job",
		Condition: condition.Condition{
			Threshold: 1,
			State:     ">",
		},
		Output: output.Output{
			Verbose:  1,
			TeamsURL: "https://mircosoft.com",
		},
		Search: search.Search{
			Host:      "https://host.com",
			Type:      "relative",
			Streamid:  "abcd12345",
			Query:     "error",
			Fields:    []string{"source", "message"},
			Frequency: 15,
		},
	}
}
