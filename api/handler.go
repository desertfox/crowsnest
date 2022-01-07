package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/jobs"
)

type NewJobReq struct {
	Name       string `json:"name"`
	QueryLink  string `json:"query"`
	OutputLink string `json:"output"`
	Threshold  int    `json:"threshold"`
}

func (s Server) createJob(r *http.Request) {
	log.Println("createJob")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var njr NewJobReq
	err = json.Unmarshal(data, &njr)
	if err != nil {
		log.Fatal(err)
	}

	job := translate(njr)

	err = jobs.AddToConfig(s.configPath, job)
	if err != nil {
		log.Fatal(err)
	}

}
