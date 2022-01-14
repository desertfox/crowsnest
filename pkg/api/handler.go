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
		State:      r.FormValue("state"),
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
	<a href="/status" target="_blank">Current Job Status Page</a>
	<form method="POST">
	<label>Job Name:</label><br />
	<input type="text" name="name"><br /><br />

	<label>GrayLog Query Link:</label><br />
	<input type="text" size="200" name="querylink" value="https://graylogquery.com?something"><br /><br />

	<label>Teams URL(<a href="https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook" target="_blank">webhook instructions</a>):</label><br />
	<input type="text" name="outputlink" value="https://teamsurl.com"><br /><br />


	<label>Verbose*:</label>
	<select name="verbose">
    	<option value="0">0</option>
		<option value="1">1</option>
	</select><br />
	<p>
		*Verbose set to 0, Output msgs will only occur when # of query results is > or < than limit and condition defined above and below this form <br />
		 Verbose set to 1, Output msgs will only occur every time Job is executed regardless of Threshold checks
	</p><br>

	<label># Of allowed query results(threshold):</label><br />
	<input type="text" name="threshold" value="0"><br /><br />
	<label>Condition*:</label>
	<select name="state">
    	<option value=">">></option>
		<option value="<"><</option>
	</select><br />
	
	<p>
		*Condition can be broken down to.<br>
		<i>let n represent # of requests from query/search results</i><br>
		<h4>Case ></h4><br>
		  n >= $Threshold<br>
		  <h4>Case <</h4><br>
		  n <= $Threshold
	</p>

	<br />
	<input type="submit" method="POST" value="Create Job">
</form>`)

	if err != nil {
		log.Fatalln(err.Error())
	}

	tmpl.Execute(w, nil)
}

func (s *Server) getStatus(w http.ResponseWriter) {
	var output template.HTML
	for _, j := range s.s.Jobs() {
		for _, tag := range j.Tags() {
			output += template.HTML(fmt.Sprintf("Tag: %v\n<br>", tag))
			output += template.HTML(
				fmt.Sprintf("<form method=\"DELETE\" action=\"/status\"><input type=\"hidden\" name=\"tag\" value=\"%v\"><input type=\"submit\" value=\"DELETE\"></form><br>", tag),
			)
		}
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
	s.event <- jobs.Event{
		Action: "reloadjobs",
		Value:  "",
	}

	s.getStatus(w)
}

func (s *Server) deleteJob(w http.ResponseWriter, r *http.Request) {
	tag := r.FormValue("tag")

	for _, j := range s.s.Jobs() {
		for _, t := range j.Tags() {
			if tag == t {
				s.event <- jobs.Event{
					Action: "DEL_TAG",
					Value:  tag,
				}

				w.Write([]byte(fmt.Sprintf("Deleted Tag %s from jobs list", tag)))

				return
			}
		}
	}
}
