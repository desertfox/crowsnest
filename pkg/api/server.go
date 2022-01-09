package api

import (
	"log"
	"net/http"

	"github.com/desertfox/crowsnest/pkg/jobs"
)

type Server struct {
	mux        *http.ServeMux
	configPath string
	jobList    *jobs.JobList
	newJob     chan jobs.Job
}

func NewServer(mux *http.ServeMux, configPath string, jobList *jobs.JobList, newJob chan jobs.Job) *Server {
	return &Server{mux, configPath, jobList, newJob}
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
		case "GET":
			log.Println("GET /")
			s.newJob <- jobs.Job{Name: "GET TEST"}
		}
	})
}
