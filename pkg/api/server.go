package api

import (
	"log"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/jobs"
	"github.com/go-co-op/gocron"
)

type Server struct {
	mux        *http.ServeMux
	newJobChan chan jobs.Job
	event      chan string
	s          *gocron.Scheduler
}

func NewServer(mux *http.ServeMux, newJobChan chan jobs.Job, event chan string, s *gocron.Scheduler) *Server {
	return &Server{mux, newJobChan, event, s}
}

func (s Server) Run() {
	s.SetupRoute()

	log.Fatal(http.ListenAndServe(":8080", s.mux))
}

func (s Server) SetupRoute() {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			s.createJob(w, r)
		case "GET":
			s.getJobForm(w)
		}
	})

	s.mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			s.getStatus(w)
		case "POST":
			s.reloadJobs(w)
		}
	})
}
