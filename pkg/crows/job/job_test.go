package job

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
		},
	}
}
