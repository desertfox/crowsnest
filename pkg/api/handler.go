package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/desertfox/crowsnest/pkg/jobs"
)

func (s *Server) createJob(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	threshold, err := strconv.Atoi(r.FormValue("threshold"))
	if err != nil {
		log.Fatal("error translating threshold to int " + err.Error())
	}

	verbose, err := strconv.Atoi(r.FormValue("verbose"))
	if err != nil {
		log.Fatal("error translating threshold to int " + err.Error())
	}

	njr := jobs.NewJobReq{
		Name:       r.FormValue("name"),
		QueryLink:  r.FormValue("querylink"),
		OutputLink: r.FormValue("outputlink"),
		Threshold:  threshold,
		Verbose:    verbose,
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
	<label>Job Name:</label><br />
	<input type="text" name="name"><br />

	<label>GrayLog Query Link:</label><br />
	<input type="text" name="querylink" value="https://graylogquery"><br />

	<label>Teams URL:</label><br />
	<input type="text" name="outputlink" value="https://teamsurl.com"><br />

	<label>Threshold for alerting:</label><br />
	<input type="text" name="threshold" value="0"><br />

	<label>Verbose*:</label><br />
	<input type="text" name="verbose" value="0"><br />
	<span>*If set to 1, will alert channel regardless of threshold</span><br />

	<br />
	<input type="submit" method="POST">
</form>`)
	if err != nil {
		log.Fatalln(err.Error())
	}

	tmpl.Execute(w, nil)
}

func (s *Server) getStatus(w http.ResponseWriter) {
	var output template.HTML
	for _, j := range s.s.Jobs() {
		output += template.HTML(fmt.Sprintf("Tags: %v\n<br>", j.Tags()))
		output += template.HTML(fmt.Sprintf("LastRun: %v<br>", j.LastRun()))
		output += template.HTML(fmt.Sprintf("NextRun: %v<br>\n", j.NextRun()))
		output += template.HTML("<br>")
	}

	tmpl, err := template.New("status_form").Parse(`
	<html>
	<h1>Crowsnest Status Now: {{ .Now}}</h1>
	{{ .Output}}
	<form method="POST">
		<label>ReloadJobs</label><br />
		<input type="submit">
	</form>	
	</html>`)
	if err != nil {
		log.Fatalln(err.Error())
	}

	tmpl.Execute(w, struct {
		Now    time.Time
		Output template.HTML
	}{
		time.Now(),
		output,
	})
}

func (s *Server) reloadJobs(w http.ResponseWriter) {
	s.event <- "reloadjobs"

	s.getStatus(w)
}
