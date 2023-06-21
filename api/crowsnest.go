package api

import (
	"encoding/json"
	"net/http"
)

func (a Api) getStatus(w http.ResponseWriter) {
	alerts := 0
	for _, v := range a.nest.Jobs() {
		if v.History.AlertCount > 0 {
			alerts++
		}
	}

	jsonStatus, err := json.Marshal(struct {
		StartTime string `json:"startTime"`
		NumJobs   int    `json:"numJobs"`
		Alerts    int    `json:"numAlerts"`
	}{
		StartTime: a.nest.StartTime.String(),
		NumJobs:   len(a.nest.Jobs()),
		Alerts:    alerts,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStatus)
}
