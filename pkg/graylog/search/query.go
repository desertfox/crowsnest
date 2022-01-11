package search

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
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
		return q.urlEncodeRelative()
	case "absolute":
		return q.urlEncodeAbsolute()
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
	params.Add("range", strconv.Itoa(q.frequnecy*60))
	params.Add("filter", fmt.Sprintf("streams:%s", q.streamid))
	params.Add("sort", "timestamp:desc")

	if len(q.fields) > 0 {
		params.Add("fields", strings.Join(q.fields, ", "))
	}

	params.Add("from", q.from.Format(time.RFC3339))
	params.Add("from", q.to.Format(time.RFC3339))
	params.Add("limit", "10000")

	return params.Encode()
}

func (q query) ExecuteSearch(authToken string) (int, error) {
	var url string

	switch q.Type {
	case "relative":
		url = fmt.Sprintf("%v/api/search/universal/relative?%v", q.host, q)
	case "absolute":
		url = fmt.Sprintf("%v/api/search/universal/absolute?%v", q.host, q)
	}

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", authToken)

	response, err := q.httpClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	//Count is probably off by one, the response is csv format iirc with the headers aka fields
	count := 0
	scanner := bufio.NewScanner(response.Body)
	scanner.Buffer([]byte{}, 1024*10048)
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func (q query) BuildSearchURL() string {
	params := url.Values{}

	params.Add("q", q.query)
	params.Add("interval", "hour")
	params.Add("rangetype", "relative")
	params.Add("relative", strconv.Itoa(q.frequnecy*60))

	return fmt.Sprintf("%s/streams/%s/search?%s", q.host, q.streamid, params.Encode())
}
