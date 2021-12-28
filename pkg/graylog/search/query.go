package search

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type query struct {
	host, query, streamid string
	frequnecy             int
	fields                []string
	httpClient            *http.Client
}

func New(host, q, streamid string, frequency int, fields []string, httpClient *http.Client) query {
	return query{host, q, streamid, frequency, fields, httpClient}
}

func (q query) urlEncode() string {
	params := url.Values{}

	params.Add("query", q.query)
	params.Add("range", strconv.Itoa(q.frequnecy*60))
	params.Add("filter", fmt.Sprintf("streams:%s", q.streamid))
	params.Add("sort", "timestamp:desc")
	params.Add("fields", strings.Join(q.fields, ", "))
	params.Add("limit", "10000")

	return params.Encode()
}

func (q query) ExecuteSearch(authToken string) (int, error) {
	url := fmt.Sprintf("%v/api/search/universal/relative?%v", q.host, q)

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
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func (q query) String() string {
	return q.urlEncode()
}

func (q query) BuildSearchURL() string {
	params := url.Values{}

	params.Add("q", q.query)
	params.Add("interval", "hour")
	params.Add("rangetype", "relative")
	params.Add("relative", strconv.Itoa(q.frequnecy*60))

	return fmt.Sprintf("%s/streams/%s/search?%s", q.host, q.streamid, params.Encode())
}
