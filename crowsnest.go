package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog"
	"github.com/desertfox/crowsnest/pkg/teams"
	"gopkg.in/yaml.v2"
)

type config struct {
	Host string `yaml:"host"`
	Jobs []job  `yaml:"jobs"`
	auth auth
}

type job struct {
	Name      string `yaml:"name"`
	Frequency int    `yaml:"frequency"`
	Option    option `yaml:"options"`
	TeamsURL  string `yaml:"teamsurl"`
	Threshold int    `yaml:"threshold"`
	Type      string `yaml:"type"`
}

type option struct {
	Streamid string   `yaml:"streamid"`
	Query    string   `yaml:"query"`
	Fields   []string `yaml:"fields"`
}

type auth struct {
	basicAuth   string
	lastUpdated time.Time
}

type reqParams struct {
	Username, Password, ConfigPath string
}

func newReqParams(u, p, c string) (reqParams, error) {
	if u == "" {
		return reqParams{}, errors.New("missing username")
	}

	if p == "" {
		return reqParams{}, errors.New("missing password")
	}

	if c == "" {
		return reqParams{}, errors.New("missing configpath")
	}

	return reqParams{u, p, c}, nil

}

func buildConfigFromENV(rp reqParams) (config, error) {
	c, err := loadConfig(rp.ConfigPath)
	if err != nil {
		return config{}, err
	}

	if err = c.InitSession(rp.Username, rp.Password); err != nil {
		return config{}, err
	}

	return c, nil
}

func loadConfig(filePath string) (config, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config{}, err
	}

	var c config
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return config{}, err
	}

	return c, nil
}

func (c *config) InitSession(u, p string) error {
	lr := graylog.NewLoginRequest(u, p, c.Host, &http.Client{})

	basicAuth, err := lr.CreateAuthHeader()
	if err != nil {
		return err
	}

	c.auth = auth{basicAuth, time.Now()}

	return nil
}

func (j job) getJob(c config) func() {
	return func() {
		fmt.Println("ExecuteJob " + j.Name)

		q := graylog.NewGLQ(c.Host, j.Name, j.Option.Query, j.Option.Streamid, c.auth.basicAuth, j.Frequency, j.Option.Fields)

		count, err := q.Execute()
		if err != nil {
			bailOut(err)
		}

		fmt.Println(time.Now(), count, q.BuildHumanURL())

		var status string
		if count >= j.Threshold {
			status = "ALERT"
		} else {
			status = "OK"
		}

		outputService := teams.BuildClient(j.TeamsURL)
		outputService.Send(j.Name, fmt.Sprintf("Status: %s\nCount: %d\nLink: [GrayLog Query](%s)\n", status, count, q.BuildHumanURL()))
	}
}

func bailOut(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
