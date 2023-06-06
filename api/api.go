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

func (a Api) Start() {
	a.jobRoutes()
	a.crowsnestRoutes()

	a.mux.Handle("/", http.FileServer(http.Dir(a.nest.HttpDir)))

	log.Fatal(http.ListenAndServe(":8080", a.mux))
}

func (a Api) jobRoutes() {
	a.mux.HandleFunc("/api/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			a.getJobs(w)
		}
	})
	a.mux.HandleFunc("/api/job", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			a.getJob(w, r)
		case "PUT":
			a.updateJob(w, r)
		case "POST":
			a.createJob(w, r)
		}
	})
}

func (a Api) crowsnestRoutes() {
	a.mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			a.getStatus(w)
		}
	})

}
