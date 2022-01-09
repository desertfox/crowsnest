package api

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/desertfox/crowsnest/pkg/jobs"
)

func (s *Server) createJob(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	threshold, err := strconv.Atoi(r.FormValue("threshold"))
	if err != nil {
		log.Fatal("error translating threshold to int " + err.Error())
	}

	njr := jobs.NewJobReq{
		Name:       r.FormValue("name"),
		QueryLink:  r.FormValue("querylink"),
		OutputLink: r.FormValue("outputlink"),
		Threshold:  threshold,
	}

	job, err := njr.TranslateToJob()
	if err != nil {
		log.Fatal("error translating job" + err.Error())
	}

	s.newJobChan <- job

	w.Write([]byte("Job Created"))
}

func (s *Server) getJobForm(w http.ResponseWriter) {
	tmpl, err := template.New("njr_form").Parse(`<h1>New Job Request Translate Form</h1>
	<form method="POST">
		<label>Name:</label><br />
		<input type="text" name="name"><br />
		<label>QueryLink:</label><br />
		<input type="text" name="querylink"><br />
		<label>OutputLink:</label><br />
		<input type="text" name="outputlink"><br />
		<label>Threshold:</label><br />
		<input type="text" name="threshold"><br />
		<input type="submit" method="POST">
	</form>`)
	if err != nil {
		log.Fatalln(err.Error())
	}

	tmpl.Execute(w, nil)
}
