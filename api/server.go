package api

import (
	"log"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/jobs"
)

type Server struct {
	mux        *http.ServeMux
	configPath string
	jobList    jobs.JobList
}

func NewServer(mux *http.ServeMux, configPath string, jobList jobs.JobList) *Server {
	return &Server{mux, configPath, jobList}
}

func (s Server) Run() {
	s.SetupRoute()

	log.Fatal(http.ListenAndServe(":8080", s.mux))
}

func (s Server) SetupRoute() {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			s.createJob(r)
		}
	})
}
