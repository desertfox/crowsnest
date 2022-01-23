package api

import (
	"log"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/crows"
)

type Api struct {
	mux  *http.ServeMux
	nest *crows.Nest
}

func New(nest *crows.Nest) Api {
	return Api{
		mux:  &http.ServeMux{},
		nest: nest,
	}
}

func (a Api) Run() {
	a.SetupRoute()

	log.Fatal(http.ListenAndServe(":8080", a.mux))
}

func (a Api) SetupRoute() {
	a.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			a.createJob(w, r)
		case "GET":
			a.getJobForm(w)
		}
	})

	a.mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			a.getStatus(w)
		case "POST":
			a.reloadJobs(w)
		}
	})

	a.mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			a.deleteJob(w, r)
		}
	})
}
