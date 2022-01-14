package jobs

import (
	"net/url"
	"strconv"
	"strings"
)

type NewJobReq struct {
	Name       string
	QueryLink  string
	OutputLink string
	Threshold  int
	Verbose    int
}

func (njr NewJobReq) TranslateToJob() (Job, error) {
	var (
		frequency  int
		typeSearch string
		fields     []string
	)

	urlObj, err := url.Parse(njr.QueryLink)
	if err != nil {
		return Job{}, err
	}

	parsedQuery := urlObj.Query()

	switch parsedQuery["rangetype"][0] {
	case "relative":
		typeSearch = "relative"
		frequency, _ = strconv.Atoi(parsedQuery["relative"][0])
		/*
			case "absolute":
				typeSearch = "absolute"
				from = parsedQuery["from"][0]
				to = parsedQuery["to"][0]
		*/
	}

	if _, ok := parsedQuery["fields"]; ok {
		fields = strings.Split(parsedQuery["fields"][0], ",")
	}

	return Job{
		njr.Name,
		Condition{
			njr.Threshold,
			"PRESENT",
		},
		Output{
			njr.Verbose,
			njr.OutputLink,
		},
		Search{
			"https://" + urlObj.Hostname(),
			typeSearch,
			getSteamId(urlObj.EscapedPath()),
			parsedQuery["q"][0],
			fields,
			frequency / 60,
		},
	}, nil
}

func getSteamId(s string) string {
	return strings.Split(s, "/")[2]
}
