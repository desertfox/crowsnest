package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/desertfox/gograylog"
)

type Search struct {
	Type     string   `yaml:"type"`
	Streamid string   `yaml:"streamid"`
	Query    string   `yaml:"query"`
	Fields   []string `yaml:"fields"`
}

type Result struct {
	Count int
	When  time.Time
}

type SearchClient interface {
	Search(gograylog.Query) ([]byte, error)
}

func (s *Search) Run(g SearchClient, frequency int) Result {
	q := gograylog.Query{
		QueryString: s.Query,
		StreamID:    s.Streamid,
		Fields:      s.Fields,
		Frequency:   frequency,
		Limit:       10000,
	}

	b, err := g.Search(q)
	if err != nil {
		log.Println(err)
		return Result{}
	}

	count := bytes.Count(b, []byte("\n"))

	if count == 0 && len(b) > 0 {
		var j map[string]interface{}
		err := json.Unmarshal(b, &j)
		if err != nil {
			fmt.Printf("Error parsing json. data:%s\nerr:%s\n", b, err)
		}

		if val, ok := j["total_results"]; ok {
			return Result{
				Count: int(val.(float64)),
				When:  time.Now(),
			}
		}
	}

	//Remove csv headers
	if count > 2 {
		count -= 1
	}

	return Result{
		Count: count,
		When:  time.Now(),
	}
}

func (s Search) BuildURL(host string, from, to time.Time) string {
	q := gograylog.Query{
		QueryString: s.Query,
		StreamID:    s.Streamid,
		Fields:      s.Fields,
		Limit:       10000,
	}

	return q.Url(host, from, to)
}

func (r Result) From(f int) time.Time {
	return r.When.Add(time.Duration(int64(-1) * int64(f) * int64(time.Minute)))
}

func (r Result) To() time.Time {
	return r.When
}
