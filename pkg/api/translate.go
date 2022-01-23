package api

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/crowsnest/pkg/crows/job/output"
	"github.com/desertfox/crowsnest/pkg/crows/job/search"
)

func translate(njr NewJobReq) (job.Job, error) {
	var (
		frequency  int
		typeSearch string
		fields     []string
	)

	urlObj, err := url.Parse(njr.QueryLink)
	if err != nil {
		return job.Job{}, err
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

	return job.Job{
		Name:      njr.Name,
		Host:      "https://" + urlObj.Hostname(),
		Frequency: frequency / 60,
		Output: output.Output{
			Verbose: njr.Verbose,
			URL:     njr.OutputLink,
		},
		Search: search.Search{
			Type:     typeSearch,
			Streamid: getSteamId(urlObj.EscapedPath()),
			Query:    parsedQuery["q"][0],
			Fields:   fields,
			Condition: search.Condition{
				Threshold: njr.Threshold,
				State:     njr.State,
			},
		},
	}, nil
}

func getSteamId(s string) string {
	return strings.Split(s, "/")[2]
}
