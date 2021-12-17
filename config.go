package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
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

func newAuth(s string) auth {
	return auth{s, time.Now()}
}

type reqParams struct {
	username, password, configPath string
}

func buildReqParams() reqParams {
	cs := "CROWSNEST_"

	return reqParams{
		username:   os.Getenv(cs + "USERNAME"),
		password:   os.Getenv(cs + "PASSWORD"),
		configPath: os.Getenv(cs + "CONFIG"),
	}
}

func buildConfigFromENV() config {
	rp := buildReqParams()

	value := reflect.ValueOf(rp)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i).Interface()
		if field == "" {
			fmt.Println("Missing ENV variable: " + "CROWSNEST_" + value.Field(i).Type().Name())
		}

	}

	c := loadConfig(rp.configPath)

	c.InitSession(rp.username, rp.password)

	return c
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

func (j job) getFrequency() int {
	return j.Frequency
}

func bailOut(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
