package graylog

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type query struct {
	host, name, query, filter, sort, limit, streamid, authToken string
	frequnecy                                                   int
	fields                                                      []string
}

func NewGLQ(host, name, q, filter, streamid, authToken string, frequency int, fields []string) query {
	return query{host, name, q, filter, "order:desc", "10000", streamid, authToken, frequency, fields}
}

func (q query) urlEncode() string {
	params := url.Values{}

	params.Add("query", q.query)
	params.Add("range", strconv.Itoa(q.frequnecy))
	params.Add("filter", q.filter)
	params.Add("sort", q.sort)
	params.Add("fields", strings.Join(q.fields, ", "))
	params.Add("limit", q.limit)

	return params.Encode()
}

func (q query) String() string {
	return q.urlEncode()
}

func (q query) Execute() (string, error) {
	url := fmt.Sprintf("%v/api/search/universal/relative?%v", q.host, q)
	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", q.authToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	return string(body), nil
}
