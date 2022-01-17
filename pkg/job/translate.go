package job

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/desertfox/crowsnest/pkg/job/condition"
	"github.com/desertfox/crowsnest/pkg/job/output"
	"github.com/desertfox/crowsnest/pkg/job/search"
)

type NewJobReq struct {
	Name       string
	QueryLink  string
	OutputLink string
	Threshold  int
	State      string
	Verbose    int
}

func (njr NewJobReq) ToJob() (Job, error) {
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
		Name: njr.Name,
		Condition: condition.Condition{
			Threshold: njr.Threshold,
			State:     njr.State,
		},
		Output: output.Output{
			Verbose:  njr.Verbose,
			TeamsURL: njr.OutputLink,
		},
		Search: search.Search{
			Host:      "https://" + urlObj.Hostname(),
			Type:      typeSearch,
			Streamid:  getSteamId(urlObj.EscapedPath()),
			Query:     parsedQuery["q"][0],
			Fields:    fields,
			Frequency: frequency / 60,
		},
	}, nil
}

func getSteamId(s string) string {
	return strings.Split(s, "/")[2]
}
