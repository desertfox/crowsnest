package search

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type query struct {
	host, query, streamid string
	frequnecy             int
	fields                []string
	Type                  string
	from                  time.Time
	to                    time.Time
	httpClient            *http.Client
}

func New(host, q, streamid string, frequency int, fields []string, t string, from, to time.Time, httpClient *http.Client) query {
	return query{host, q, streamid, frequency, fields, t, from, to, httpClient}
}

func (q query) String() string {
	switch q.Type {
	case "relative":
		if os.Getenv("CROWSNEST_DEBUG") == "1" {
			log.Printf("%v/api/search/universal/relative?%v", q.host, q.urlEncodeRelative())
		}
		return fmt.Sprintf("%v/api/search/universal/relative?%v", q.host, q.urlEncodeRelative())
	case "absolute":
		if os.Getenv("CROWSNEST_DEBUG") == "1" {
			log.Printf("%v/api/search/universal/absolute?%v", q.host, q.urlEncodeAbsolute())
		}
		return fmt.Sprintf("%v/api/search/universal/absolute?%v", q.host, q.urlEncodeAbsolute())
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

func (q query) urlEncodeAbsolute() string {
	params := url.Values{}

	params.Add("query", q.query)
	params.Add("filter", fmt.Sprintf("streams:%s", q.streamid))
	params.Add("sort", "timestamp:desc")

	if len(q.fields) > 0 {
		params.Add("fields", strings.Join(q.fields, ", "))
	}

	params.Add("from", q.from.Format(time.RFC3339))
	params.Add("to", q.to.Format(time.RFC3339))
	params.Add("limit", "10000")

	return params.Encode()
}

func (q query) ExecuteSearch(authToken string) (int, error) {
	request, _ := http.NewRequest("GET", q.String(), nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", authToken)

	response, err := q.httpClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("CROWSNEST_DEBUG") == "log" {
		log.Println(string(body))
	}

	return strings.Count(string(body), "\n"), nil
}

func (q query) BuildSearchURL() string {
	params := url.Values{}

	params.Add("q", q.query)

	if q.Type == "relative" {
		params.Add("interval", "minute")
		params.Add("rangetype", "relative")
		params.Add("relative", strconv.Itoa(q.frequnecy*60))
	}

	if q.Type == "absolute" {
		params.Add("rangetype", "absolute")
		params.Add("from", q.from.Format(time.RFC3339))
		params.Add("to", q.to.Format(time.RFC3339))
	}

	return fmt.Sprintf("%s/streams/%s/search?%s", q.host, q.streamid, params.Encode())
}
