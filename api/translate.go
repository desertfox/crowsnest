package api

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/desertfox/crowsnest/pkg/crows/job"
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
		frequency, err = strconv.Atoi(parsedQuery["relative"][0])
		if err != nil {
			return job.Job{}, err
		}
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

	streamid, err := getSteamId(urlObj.EscapedPath())
	if err != nil {
		return job.Job{}, err
	}

	return job.Job{
		Name:      njr.Name,
		Host:      "https://" + urlObj.Hostname(),
		Frequency: frequency / 60,
		Verbose:   njr.Verbose,
		Teams: job.Teams{
			Url:  njr.OutputLink,
			Name: njr.OutputName,
		},
		Search: job.Search{
			Type:     typeSearch,
			Streamid: streamid,
			Query:    parsedQuery["q"][0],
			Fields:   fields,
		},
		Condition: job.Condition{
			Threshold: njr.Threshold,
			State:     njr.State,
		},
	}, nil
}

func getSteamId(s string) (string, error) {
	parts := strings.Split(s, "/")
	if len(parts) == 1 {
		return "", errors.New("unexpected format: " + s)
	}
	return strings.Split(s, "/")[2], nil
}
