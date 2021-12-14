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
	host string `yaml:"host"`
	jobs []job  `yaml:"jobs"`
	auth auth
}

type job struct {
	name      string `yaml:"name"`
	frequency int    `yaml:"frequency"`
	option    option
}

type option struct {
	steamid string   `yaml:"steamid"`
	query   string   `yaml:"query"`
	fields  []string `yaml:"fields"`
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
	bailOut(err)

	var c config
	err = yaml.Unmarshal(file, &c)
	bailOut(err)

	return c
}

func (c *config) InitSession(u, p string) {
	lr := graylog.NewLoginRequest(u, p, c.host, &http.Client{})

	basicAuth, err := lr.CreateAuthHeader()
	if err != nil {
		bailOut(err)
	}

	c.auth = newAuth(basicAuth)
}

func (j job) getFunc(c config) func() {
	return func() {
		fmt.Println("ExecuteJob " + j.name)

		q := graylog.NewGLQ(c.host, j.name, j.option.query, "", j.option.steamid, c.auth.basicAuth, j.frequency, j.option.fields)

		q.Execute()
	}
}

func (j job) getFrequency() int {
	return j.frequency
}
