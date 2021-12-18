package main

import "time"

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

type LoginRequest interface {
	CreateAuthHeader() (string, error)
}

type Query interface {
	Execute() (int, error)
	BuildHumanURL() string
}

type Report interface {
	Send(string, string) error
}
