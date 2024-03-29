package crows

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/desertfox/crowsnest/pkg/crows/cron"
	"github.com/desertfox/crowsnest/pkg/crows/job"
	"github.com/desertfox/gograylog"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	graylog  *graylog
	JobsPath string `yaml:"jobspath"`
	TeamsURL string `yaml:"teamsurl"`
	HttpDir  string `yaml:"httpdir"`
}

type graylog struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (c *Config) Load(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return fmt.Errorf("unable to read Config, %w", err)
	}

	err = yaml.Unmarshal(file, &c.graylog)
	if err != nil {
		return fmt.Errorf("unable to read graylog config, %w", err)
	}

	return nil
}

func (c *Config) BuildNest() *Nest {
	z, _ := zap.NewDevelopment()

	loc, _ := time.LoadLocation("America/Chicago")

	return &Nest{
		Config: c,
		list:   c.BuildList(),
		schedule: &cron.Schedule{
			Scheduler: gocron.NewScheduler(time.UTC),
		},
		MSTeamsClient: c.buildTeamsClient(),
		GrayLogClient: c.buildGraylogClient(),
		log:           z.Sugar(),
		StartTime:     time.Now().In(loc),
	}
}

func (c *Config) buildGraylogClient() *gograylog.Client {
	g := &gograylog.Client{
		Host:     c.graylog.Host,
		Username: c.graylog.Username,
		Password: c.graylog.Password,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	return g
}

func (c *Config) buildTeamsClient() *goteamsnotify.TeamsClient {
	return goteamsnotify.NewTeamsClient()
}

func (c *Config) BuildList() *job.List {
	return &job.List{
		File: c.JobsPath,
		Jobs: make([]*job.Job, 0),
	}
}
