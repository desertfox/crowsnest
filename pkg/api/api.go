package api

import (
	"log"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/job"
)

type Api struct {
	Mux       *http.ServeMux
	Scheduler *job.Scheduler
}

func New(j *job.Scheduler) Api {
	return Api{
		Mux:       &http.ServeMux{},
		Scheduler: j,
	}
}

func (a Api) Run() {
	a.SetupRoute()

	log.Fatal(http.ListenAndServe(":8080", a.Mux))
}

func (a Api) SetupRoute() {
	a.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			a.createJob(w, r)
		case "GET":
			a.getJobForm(w)
		case "DELETE":
			a.deleteJob(w, r)
		}
	})

	a.Mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			a.getStatus(w)
		case "POST":
			a.reloadJobs(w)
		}
	})

	a.Mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			a.deleteJob(w, r)
		}
	})
}
