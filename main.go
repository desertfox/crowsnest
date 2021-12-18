package main

import (
	"net/http"
	"os"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog"
	"github.com/desertfox/crowsnest/pkg/teams"
	"github.com/go-co-op/gocron"
)

func main() {
	rp, err := newReqParams(os.Getenv("CROWSNEST_USERNAME"), os.Getenv("CROWSNEST_PASSWORD"), os.Getenv("CROWSNEST_CONFIG"))
	if err != nil {
		bailOut(err)
	}

	c, err := buildConfig(rp)
	if err != nil {
		bailOut(err)
	}

	lr := graylog.NewLoginRequest(rp.Username, rp.Password, c.Host, &http.Client{})
	if err = c.buildSession(lr); err != nil {
		bailOut(err)
	}

	s := gocron.NewScheduler(time.UTC)

	for _, j := range c.Jobs {
		q := graylog.NewGLQ(c.Host, j.Name, j.Option.Query, j.Option.Streamid, c.auth.basicAuth, j.Frequency, j.Option.Fields)
		r := teams.BuildClient(j.TeamsURL)

		s.Every(j.Frequency).Minutes().Do(j.getJob(q, r))
	}

	s.StartBlocking()
}
