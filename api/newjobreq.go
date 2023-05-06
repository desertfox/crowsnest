package api

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/desertfox/crowsnest/pkg/crows/job"
)

type NewJobReq struct {
	//Job Name
	Name string `json:"name"`
	//Graylog link
	GraylogLink string `json:"graylogLink"`
	//Teams URL
	TeamsUrl string `json:"teamsUrl"`
	//Teams Room Name
	TeamsRoomName string `json:"teamsRoomName"`
	//Integer value used to determine boundary
	Threshold int `json:"threshold"`
	//Operator used to determine boundary
	Operator string `json:"operator"`
	//Determines if all runs should send msg or just alerts
	Verbose bool `json:"verbose"`
	//start time off set, 24:00
	OffSet string `json:"offset"`
}

func translate(njr NewJobReq) (job.Job, error) {
	var (
		frequency int
		fields    []string
	)

	urlObj, err := url.Parse(njr.GraylogLink)
	if err != nil {
		return job.Job{}, err
	}

	parsedQuery := urlObj.Query()

	fmt.Printf("%#v", parsedQuery)

	switch parsedQuery["rangetype"][0] {
	case "relative":
		frequency, err = strconv.Atoi(parsedQuery["relative"][0])
		if err != nil {
			return job.Job{}, err
		}
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
			Url:  njr.TeamsUrl,
			Name: njr.TeamsRoomName,
		},
		Search: job.Search{
			Streamid: streamid,
			Query:    parsedQuery["q"][0],
			Fields:   fields,
		},
		Condition: job.Condition{
			Threshold: njr.Threshold,
			Operator:  njr.Operator,
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
