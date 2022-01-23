package api

import (
	"testing"

	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/crowsnest/pkg/crows/job/output"
	"github.com/desertfox/crowsnest/pkg/crows/job/search"
)

var (
	name              string = "Test Name"
	frequency         int    = 1440
	relativeQueryLink string = "https://desertfox.dev/streams/5555555/search?rangetype=relative&fields=source%2Cmessage&width=1920&highlightMessage=&relative=86400&q=%22Error+Checking+Out%22"
	outputLink        string = "https://not.desertfox.dev"
	threshold         int    = 5

	host                        = "https://desertfox.dev"
	searchTypeRelative string   = "relative"
	streamid           string   = "5555555"
	q                  string   = "\"Error Checking Out\""
	fields             []string = []string{"source", "message"}
	verbose            int      = 1
)

func Test_translate(t *testing.T) {

	t.Run("translate", func(t *testing.T) {
		njr := NewJobReq{name, relativeQueryLink, outputLink, threshold, "<", verbose}
		got, gotErr := translate(njr)
		if gotErr != nil {
			t.Errorf("error got: %#v", gotErr)
		}

		want := job.Job{
			Name:      name,
			Frequency: frequency,
			Host:      host,
			Output: output.Output{
				URL:     outputLink,
				Verbose: verbose,
			},
			Search: search.Search{
				Type:     searchTypeRelative,
				Streamid: streamid,
				Query:    q,
				Fields:   fields,
				Condition: search.Condition{
					Threshold: threshold,
				},
			},
		}

		if got.Name != want.Name {
			t.Errorf("error got: %#v, want %#v", got.Name, want.Name)
		}

		if got.Frequency != want.Frequency {
			t.Errorf("error got: %#v, want %#v", got.Frequency, want.Frequency)
		}

		if got.Search.Condition.Threshold != want.Search.Condition.Threshold {
			t.Errorf("error got: %#v, want %#v", got.Search.Condition.Threshold, want.Search.Condition.Threshold)
		}

		if got.Output.URL != want.Output.URL {
			t.Errorf("error got: %#v, want %#v", got.Output.URL, want.Output.URL)
		}

		if got.Host != want.Host {
			t.Errorf("error got: %#v, want %#v", got.Host, want.Host)
		}

		if got.Search.Type != want.Search.Type {
			t.Errorf("error got: %#v, want %#v", got.Search.Type, want.Search.Type)
		}

		if got.Search.Streamid != want.Search.Streamid {
			t.Errorf("error got: %#v, want %#v", got.Search.Streamid, want.Search.Streamid)
		}

		if got.Search.Query != want.Search.Query {
			t.Errorf("error got: %#v, want %#v", got.Search.Query, want.Search.Query)
		}

		if len(got.Search.Fields) != len(want.Search.Fields) {
			t.Errorf("error got: %#v, want %#v", len(got.Search.Fields), len(want.Search.Fields))
		}

		for index, got := range got.Search.Fields {
			if got != want.Search.Fields[index] {
				t.Errorf("error got: %#v, want %#v", got, want.Search.Fields[index])
			}
		}

	})
}
