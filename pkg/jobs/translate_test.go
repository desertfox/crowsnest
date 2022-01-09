package jobs

import (
	"testing"
)

var (
	name              string = "Test Name"
	frequency         int    = 86400
	relativeQueryLink string = "https://desertfox.dev/streams/5555555/search?rangetype=relative&fields=source%2Cmessage&width=1920&highlightMessage=&relative=86400&q=%22Error+Checking+Out%22"
	outputLink        string = "https://not.desertfox.dev"
	threshold         int    = 5

	host                        = "https://desertfox.dev"
	searchTypeRelative string   = "relative"
	streamid           string   = "5555555"
	query              string   = "\"Error Checking Out\""
	fields             []string = []string{"source", "message"}
)

func Test_translate(t *testing.T) {

	t.Run("translate", func(t *testing.T) {

		njr := NewJobReq{name, relativeQueryLink, outputLink, threshold}
		got, gotErr := njr.TranslateToJob()
		if gotErr != nil {
			t.Errorf("error got: %#v", gotErr)
		}

		want := Job{
			Name:      name,
			Frequency: frequency,
			Threshold: threshold,
			TeamsURL:  outputLink,
			Search: SearchOptions{
				Host:     host,
				Type:     searchTypeRelative,
				Streamid: streamid,
				Query:    query,
				Fields:   fields,
			},
		}

		if got.Name != want.Name {
			t.Errorf("error got: %#v, want %#v", got.Name, want.Name)
		}

		if got.Frequency != want.Frequency {
			t.Errorf("error got: %#v, want %#v", got.Frequency, want.Frequency)
		}

		if got.Threshold != want.Threshold {
			t.Errorf("error got: %#v, want %#v", got.Threshold, want.Threshold)
		}

		if got.TeamsURL != want.TeamsURL {
			t.Errorf("error got: %#v, want %#v", got.TeamsURL, want.TeamsURL)
		}

		if got.Search.Host != want.Search.Host {
			t.Errorf("error got: %#v, want %#v", got.Search.Host, want.Search.Host)
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
