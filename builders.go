package main

import (
	"errors"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	errReqParams = errors.New("missing param")
)

func newReqParams(u, p, c string) (reqParams, error) {
	for _, v := range []string{u, p, c} {
		if v == "" {
			return reqParams{}, errReqParams
		}
	}

	return reqParams{u, p, c}, nil

}

func buildConfig(rp reqParams) (config, error) {
	file, err := ioutil.ReadFile(rp.ConfigPath)
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

func (c *config) buildSession(lr LoginRequest) error {
	basicAuth, err := lr.CreateAuthHeader()
	if err != nil {
		return err
	}

	c.auth = auth{basicAuth, time.Now()}

	return nil
}
