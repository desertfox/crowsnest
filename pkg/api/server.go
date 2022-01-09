package api

import (
	"log"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/jobs"
)

type Server struct {
	mux        *http.ServeMux
	newJobChan chan jobs.Job
}

func NewServer(mux *http.ServeMux, newJobChan chan jobs.Job) *Server {
	return &Server{mux, newJobChan}
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
}
