package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/desertfox/crowsnest/pkg/crows/job"
)

type NewJobReq struct {
	Name       string
	QueryLink  string
	OutputLink string
	OutputName string
	Threshold  int
	State      string
	Verbose    int
	OffSet     string
}

func (a Api) createJob(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	threshold, err := strconv.Atoi(r.FormValue("threshold"))
	if err != nil {
		w.Write([]byte("error translating threshold to int " + err.Error()))
		return
	}

	verbose, err := strconv.Atoi(r.FormValue("verbose"))
	if err != nil {
		w.Write([]byte("error translating verbose to int " + err.Error()))
		return
	}

	var (
		outputlink    string
		outputname    string
		existRoomName string = r.FormValue("exist")
	)

	if existRoomName != "" {
		for _, j := range a.nest.Jobs() {
			if j.Output.Teams.Name == existRoomName {
				outputlink = j.Output.Teams.Url
				outputname = j.Output.Teams.Name
				break
			}
		}
	} else {
		outputlink = r.FormValue("outputlink")
		outputname = r.FormValue("outputname")
	}

	njr := NewJobReq{
		Name:       r.FormValue("name"),
		QueryLink:  r.FormValue("querylink"),
		OutputLink: outputlink,
		OutputName: outputname,
		Threshold:  threshold,
		State:      r.FormValue("state"),
		Verbose:    verbose,
		OffSet:     r.FormValue("offset"),
	}

	j, err := translate(njr)
	if err != nil {
		w.Write([]byte("error translating job" + err.Error()))
		return
	}

	a.nest.HandleEvent(job.Event{
		Action: job.Add,
		Job:    &j,
	})

	w.Write([]byte("Job Created"))
}

func (a Api) getJobForm(w http.ResponseWriter) {
	tmpl, err := template.New("njr_form").Parse(`
	<html>
	<h1>New Job Request Translate Form</h1>
	<a href="/status" target="_blank">Current Job Status Page</a>
	<div class="cn-form">
	<form method="POST">
	<label>Job Name:</label>
	<input type="text" name="name"><br /><br />
	<label>GrayLog Query Link:</label><br />
	<textarea id="querylink" name="querylink" value="https://graylogquery.com?something"></textarea><br /><br />
	<label>OffSet ##:##:</label>
	<input type="text" name="offset"><br /><br />
	<label>Teams Room Name/Label:</label>
	<input type="text" name="outputname"> <br>
	<label>Teams URL:</label>
	<input type="text" name="outputlink" value="https://teamsurl.com">
	<a href="https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook" target="_blank">webhook instructions</a>
	<br>
	OR<br>
	<select name="exist">
		<option value="" selected>None</option>
		{{range .RoomNames}}
			<option value="{{ . }}">{{ . }}</option>
		{{end}}
	</select>
	<br>
	<br>
	<label>When to message room:</label>
	<select name="verbose">
    	<option value="0">outside of threshold</option>
		<option value="1">all</option>
	</select>
	<br><br>
	<label>Number of cases to alert(n):</label>
	<input type="text" name="threshold" value="0"><br /><br />
	<label>Condition:</label>
	<select name="state">
    	<option value=">">Alert when cases >= n#</option>
		<option value="<">Alert when cases <= n#</option>
	</select><br />
	<br />
	<input type="submit" method="POST" value="Create Job">
</form>
</div>
</html>`)

	if err != nil {
		log.Fatalln(err.Error())
	}

	var (
		roomNames []string
		uniqRooms map[string]bool = make(map[string]bool)
	)
	for _, j := range a.nest.Jobs() {
		if _, ok := uniqRooms[j.Output.Teams.Name]; !ok {
			uniqRooms[j.Output.Teams.Name] = true
			roomNames = append(roomNames, j.Output.Teams.Name)
		}
	}

	tmpl.Execute(w, struct {
		RoomNames []string
	}{
		roomNames,
	})
}

func (a Api) getStatus(w http.ResponseWriter) {
	central, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatalf(err.Error())
	}

	var output template.HTML
	for _, j := range a.nest.Jobs() {
		var results template.HTML = template.HTML(fmt.Sprintf("Average: %d<br>", j.History.Avg()))
		for i, r := range j.History.Results() {
			results += template.HTML(fmt.Sprintf(`Index: %d, When: %s, Count: %d, Link: <a href="%s" target="_blank">GrayLog</a><br>`,
				i,
				r.When.In(central).Format(time.RFC822),
				r.Count,
				j.Search.BuildURL(r.From(j.Frequency), r.To()),
			))
		}

		output += template.HTML(fmt.Sprintf(`
				<div style="border-style: solid">
				<label>Name: %s</label><br>
				<label>Frequency: %d min(s)</label><br>
				<label>Threshold: %d </label><br>
				<label>Operator: %s </label><br>
				<label>LastRun: %s</label><br>
				<label>NextRun: %s</label><br>
				<label>Results:<br>
				%s
				</label><br>
				<form method="POST" action="/delete">
					<input type="hidden" name="name" value="%s">
					<input type="submit" value="DELETE">
				</form></div><br>`,
			j.Name,
			j.Frequency,
			j.Condition.Threshold,
			j.Condition.State,
			a.nest.LastRun(j.Name).In(central).Format(time.RFC822),
			a.nest.NextRun(j.Name).In(central).Format(time.RFC822),
			results,
			j.Name,
		))
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
		Now    string
		Output template.HTML
	}{
		time.Now().Format(time.RFC822),
		output,
	})
}

func (a Api) reloadJobs(w http.ResponseWriter) {
	a.nest.HandleEvent(job.Event{
		Action: job.Reload,
	})

	a.getStatus(w)
}

func (a Api) deleteJob(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	tmpl, err := template.New("del_job").Parse(`
	<html>
	<h1>Crowsnest Job Removed</h1>
	{{ .Output}}
	</html>`)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, j := range a.nest.Jobs() {
		if name == j.Name {
			a.nest.HandleEvent(job.Event{
				Action: job.Del,
				Job:    &job.Job{Name: name},
			})

			output := template.HTML(
				fmt.Sprintf(`Deleted Tag %s from jobs list<br>
				<a href="/status">Job Status</a><a href="/">Job Form</a>`, name),
			)

			tmpl.Execute(w, struct {
				Output template.HTML
			}{
				output,
			})
			return
		}
	}

	log.Printf("Job not found in scheduled jobs list: %s", name)
}
