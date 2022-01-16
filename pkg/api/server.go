package api

import (
	"log"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/jobs"
	"github.com/go-co-op/gocron"
)

type Server struct {
	Mux       *http.ServeMux
	Event     chan jobs.Event
	Scheduler *gocron.Scheduler
}

func (s Server) Run() {
	s.SetupRoute()

	log.Fatal(http.ListenAndServe(":8080", s.Mux))
}

func (s Server) SetupRoute() {
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			s.createJob(w, r)
		case "GET":
			s.getJobForm(w)
		case "DELETE":
			s.deleteJob(w, r)
		}
	})

	s.Mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			s.getStatus(w)
		case "POST":
			s.reloadJobs(w)
		}
	})

	s.Mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			s.deleteJob(w, r)
		}
	})
}
