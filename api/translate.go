package api

import (
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/desertfox/crowsnest/pkg/jobs"
)

func translate(njr NewJobReq) jobs.Job {
	var (
		frequency            int
		typeSearch, from, to string
		fields               []string
	)

	urlObj, err := url.Parse(njr.QueryLink)
	if err != nil {
		log.Fatal(err)
	}

	parsedQuery := urlObj.Query()

	switch parsedQuery["rangetype"][0] {
	case "relative":
		typeSearch = "relative"
		frequency, _ = strconv.Atoi(parsedQuery["relative"][0])
	case "absolute":
		typeSearch = "absolute"
		from = parsedQuery["from"][0]
		to = parsedQuery["to"][0]
	}

	if _, ok := parsedQuery["fields"]; ok {
		fields = strings.Split(parsedQuery["fields"][0], ",")
	}

	so := jobs.NewSearchOptions("https://"+urlObj.Hostname(), typeSearch, getSteamId(urlObj.EscapedPath()), parsedQuery["q"][0], fields, from, to)

	return jobs.NewJob(njr.Name, frequency, njr.Threshold, njr.OutputLink, so)
}

func getSteamId(s string) string {
	return strings.Split(s, "/")[2]
}
