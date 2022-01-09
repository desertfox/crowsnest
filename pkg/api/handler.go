package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/jobs"
)

func (s *Server) createJob(r *http.Request) {
	log.Println("createJob")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var njr jobs.NewJobReq
	err = json.Unmarshal(data, &njr)
	if err != nil {
		log.Fatal(err)
	}

	job, err := njr.TranslateToJob()
	if err != nil {
		log.Fatal(err)
	}

	err = s.jobList.Add(job)
	if err != nil {
		log.Fatal(err)
	}

	err = s.jobList.WriteConfig(s.configPath)
	if err != nil {
		log.Fatal(err)
	}
}
