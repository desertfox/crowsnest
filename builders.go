package main

import (
	"errors"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

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
