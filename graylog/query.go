package graylog

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	grayLogDateFormat   string = "2006-01-02T15:04:05.000Z"
	relativeStrTempalte string = "%v/api/search/universal/relative?%v"
	//absoluteStrTempalte string = "%v/api/search/universal/absolute?%v"
)

type Query struct {
	Host, Query, Streamid, Type string
	Frequency                   int
	Fields                      []string
}

func (q Query) String() string {
	switch q.Type {
	case "relative":
		return fmt.Sprintf(relativeStrTempalte, q.Host, q.urlEncodeRelative())
	}
	return ""
}

func (q Query) urlEncodeRelative() string {
	params := url.Values{}

	params.Add("query", q.Query)
	params.Add("range", strconv.Itoa(q.Frequency*60))
	params.Add("filter", fmt.Sprintf("streams:%s", q.Streamid))
	params.Add("sort", "timestamp:desc")
	params.Add("fields", strings.Join(q.Fields, ", "))
	params.Add("limit", "10000")

	return params.Encode()
}

func (q Query) execute(authToken string, httpClient *http.Client) ([]byte, error) {
	request, _ := http.NewRequest("GET", q.String(), nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", authToken)

	response, err := httpClient.Do(request)
	if err != nil {
		log.Printf("error making request:%#v, error:%v", response, err)
		return []byte{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("error reading response:%#v, error:%v", response, err)
		return []byte{}, err
	}

	return body, nil
}

func (q Query) toURL(from, to time.Time) string {
	params := url.Values{}

	params.Add("q", q.Query)
	params.Add("fields", strings.Join(q.Fields, ", "))

	params.Add("rangetype", "absolute")
	params.Add("from", from.Format(grayLogDateFormat))
	params.Add("to", to.Format(grayLogDateFormat))

	//params.Add("interval", "minute")
	//params.Add("rangetype", "relative")
	//params.Add("relative", strconv.Itoa(q.Frequency*60))

	return fmt.Sprintf("%s/streams/%s/search?%s", q.Host, q.Streamid, params.Encode())
}
