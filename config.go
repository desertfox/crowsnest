package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/desertfox/crowsnest/pkg/graylog"
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
	option    option
}

type option struct {
	Steamid string   `yaml:"steamid"`
	Query   string   `yaml:"query"`
	Fields  []string `yaml:"fields"`
}

type auth struct {
	basicAuth   string
	lastUpdated time.Time
}

func newAuth(s string) auth {
	return auth{s, time.Now()}
}

func loadConfig(filePath string) config {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		bailOut(err)
	}

	var c config
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		bailOut(err)
	}

	return c
}

func (c *config) InitSession(u, p string) {
	lr := graylog.NewLoginRequest(u, p, c.Host, &http.Client{})

	basicAuth, err := lr.CreateAuthHeader()
	if err != nil {
		bailOut(err)
	}

	c.auth = newAuth(basicAuth)
}

func (j job) getFunc(c config) func() {
	return func() {
		fmt.Println("ExecuteJob " + j.Name)

		q := graylog.NewGLQ(c.Host, j.Name, j.option.Query, "", j.option.Steamid, c.auth.basicAuth, j.Frequency, j.option.Fields)

		q.Execute()
	}
}

func (j job) getFrequency() int {
	return j.Frequency
}
