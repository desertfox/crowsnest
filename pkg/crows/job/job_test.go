package job

func testJob() Job {
	return Job{
		Name: "Test Job",
		Condition: Condition{
			Threshold: 1,
			State:     ">",
		},
		Output: Output{
			Verbose:  1,
			TeamsURL: "https://mircosoft.com",
		},
		Search: Search{
			Host:      "https://host.com",
			Type:      "relative",
			Streamid:  "abcd12345",
			Query:     "error",
			Fields:    []string{"source", "message"},
			Frequency: 15,
		},
	}
}
