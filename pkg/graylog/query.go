package graylog

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type query struct {
	host, name, query, streamid string
	frequnecy                   int
	fields                      []string
}

func (j job) newQuery(host string) query {
	return query{host, j.Name, j.Option.Query, j.Option.Streamid, j.Frequency, j.Option.Fields}
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

func (q query) String() string {
	return q.urlEncode()
}

func (q query) Execute(authToken string) (int, error) {
	url := fmt.Sprintf("%v/api/search/universal/relative?%v", q.host, q)
	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", authToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

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

func (q query) BuildHumanURL() string {
	params := url.Values{}

	params.Add("q", q.query)
	params.Add("interval", "hour")
	params.Add("rangetype", "relative")
	params.Add("relative", strconv.Itoa(q.frequnecy*60))

	return q.host + "/streams/" + q.streamid + "/search?" + params.Encode()
}
