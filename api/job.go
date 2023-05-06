package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/desertfox/crowsnest/pkg/crows"
)

func (a Api) getJobs(w http.ResponseWriter) {
	jsonJobs, err := json.Marshal(a.nest.Jobs())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonJobs)
}

func (a Api) getJob(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	jobIdStr := parts[len(parts)-1]
	jobId, _ := strconv.Atoi(jobIdStr)

	if jobId >= len(a.nest.Jobs())+1 {
		http.Error(w, fmt.Sprintf("Job index not found %s %d", jobIdStr, len(a.nest.Jobs())), http.StatusInternalServerError)
		return
	}

	job := a.nest.Jobs()[jobId]

	jsonJobs, err := json.Marshal(job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonJobs)
}

func (a Api) updateJob(w http.ResponseWriter, r *http.Request) {
	/*
		parts := strings.Split(r.URL.Path, "/")
		jobIdStr := parts[len(parts)-1]
		jobId, _ := strconv.Atoi(jobIdStr)
	*/

	jsonJobs, err := json.Marshal(a.nest.Jobs())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonJobs)
}

func (a Api) createJob(w http.ResponseWriter, r *http.Request) {
	njr := &NewJobReq{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, njr)

	j, err := translate(*njr)
	if err != nil {
		http.Error(w, "error translating job"+err.Error(), http.StatusBadRequest)
		return
	}

	a.nest.HandleEvent(crows.Event{
		Action: crows.Add,
		Job:    &j,
	})

	b, _ := json.Marshal(struct {
		Id int `json:"id"`
	}{
		Id: len(a.nest.Jobs()) - 1,
	})
	w.Write(b)
}
