package api

import (
	"encoding/json"
	"net/http"
)

func (a Api) getStatus(w http.ResponseWriter) {
	jsonStatus, err := json.Marshal(struct {
		StartTime string `json:"startTime"`
		NumJobs   int    `json:"numJobs"`
	}{
		StartTime: a.nest.StartTime.String(),
		NumJobs:   len(a.nest.Jobs()),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStatus)
}
