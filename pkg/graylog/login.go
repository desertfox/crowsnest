package graylog

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const sessionsPath string = "api/system/sessions"

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	client   *http.Client
}

func NewLoginRequest(u, p, h string, c *http.Client) loginRequest {
	return loginRequest{u, p, h, c}
}

func (lr loginRequest) getSessionId() (string, error) {
	jsonData, err := json.Marshal(lr)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%v/%v", lr.Host, sessionsPath)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := lr.client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	var data map[string]string
	_ = json.Unmarshal(body, &data)

	return data["session_id"], nil
}

func (lr loginRequest) CreateAuthHeader() (string, error) {
	sessionId, err := lr.getSessionId()
	if err != nil {
		return "", err
	}

	return "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:session", sessionId))), nil
}
