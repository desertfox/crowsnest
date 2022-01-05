package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_createJob(t *testing.T) {

	t.Run("simple request", func(t *testing.T) {
		njr := NewJobReq{"test", "http", "httpo", 5}
		data, _ := json.Marshal(&njr)
		r := ioutil.NopCloser(bytes.NewReader(data))
		test := &http.Request{
			Body: r,
		}
		s := Server{}
		s.createJob(test)
	})
}
