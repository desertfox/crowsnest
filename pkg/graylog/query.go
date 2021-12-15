package graylog

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type query struct {
	host, name, query, streamid, authToken string
	frequnecy                              int
	fields                                 []string
}

func NewGLQ(host, name, q, streamid, authToken string, frequency int, fields []string) query {
	return query{host, name, q, streamid, authToken, frequency, fields}
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

func (q query) Execute() (int, error) {
	url := fmt.Sprintf("%v/api/search/universal/relative?%v", q.host, q)
	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", q.authToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := response.Body.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
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
