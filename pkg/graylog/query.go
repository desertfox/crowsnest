package graylog

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	//grayLogDateFormat   string = "2006-02-02T15:04:05.000Z"
	relativeStrTempalte string = "%v/api/search/universal/relative?%v"
	//absoluteStrTempalte string = "%v/api/search/universal/absolute?%v"
)

type query struct {
	host, query, streamid string
	frequnecy             int
	fields                []string
	Type                  string
}

func (q query) String() string {
	switch q.Type {
	case "relative":
		if os.Getenv("CROWSNEST_DEBUG") == "1" {
			log.Printf(relativeStrTempalte, q.host, q.urlEncodeRelative())
		}
		return fmt.Sprintf(relativeStrTempalte, q.host, q.urlEncodeRelative())
	}
	return ""
}

func (q query) urlEncodeRelative() string {
	params := url.Values{}

	params.Add("query", q.query)
	params.Add("range", strconv.Itoa(q.frequnecy*60))
	params.Add("filter", fmt.Sprintf("streams:%s", q.streamid))
	params.Add("sort", "timestamp:desc")
	params.Add("fields", strings.Join(q.fields, ", "))
	params.Add("limit", "10000")

	return params.Encode()
}

func (q query) execute(authToken string, httpClient *http.Client) (int, error) {
	request, _ := http.NewRequest("GET", q.String(), nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", authToken)

	response, err := httpClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("CROWSNEST_DEBUG") == "log" {
		log.Println("Body: " + string(body))
	}

	count := strings.Count(string(body), "\n")

	//Results include CSV header row, -1 to correct count.
	if count > 1 {
		count -= 1
	}

	return count, nil
}

func (q query) toURL() string {
	params := url.Values{}

	params.Add("q", q.query)
	params.Add("interval", "minute")
	params.Add("rangetype", "relative")
	params.Add("relative", strconv.Itoa(q.frequnecy*60))

	return fmt.Sprintf("%s/streams/%s/search?%s", q.host, q.streamid, params.Encode())
}
