package graylog

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const sessionsPath string = "api/system/sessions"

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

type auth struct {
	basicAuth string
	updated   time.Time
}

func (c *Client) getSessionHeader() string {
	//check if token is old
	if 1 == 0 {
		sessionId, err := c.sessionIdRequest()
		if err != nil {
			panic(err.Error())
		}

		c.auth.basicAuth = c.createAuthHeader(sessionId)
		c.auth.updated = time.Now()
	}
	// Token is good
	return c.auth.basicAuth
}

func (c Client) sessionIdRequest() (string, error) {
	jsonData, err := json.Marshal(c.lr)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%v/%v", c.lr.Host, sessionsPath)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	var data map[string]string
	_ = json.Unmarshal(body, &data)

	return data["session_id"], nil
}

func (c Client) createAuthHeader(sessionId string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:session", sessionId)))
}
