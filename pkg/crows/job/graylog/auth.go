package graylog

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const sessionsPath string = "api/system/sessions"

type session struct {
	basicAuth    string
	updated      time.Time
	loginRequest *loginRequest
}

type loginRequest struct {
	Host       string `json:"host"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	httpClient *http.Client
}

var (
	lock               = &sync.Mutex{}
	sessionInstanceMap = make(map[string]*session)
)

func newSession(h, u, p string, httpClient *http.Client) *session {
	lock.Lock()
	defer lock.Unlock()
	if _, exists := sessionInstanceMap[h]; !exists {
		sessionInstanceMap[h] = &session{"", time.Now(), &loginRequest{h, u, p, httpClient}}
	}

	return sessionInstanceMap[h]
}

func (s *session) authHeader() string {
	sessionId, err := s.loginRequest.execute()
	if err != nil {
		panic(err.Error())
	}

	s.basicAuth = createAuthHeader(sessionId)
	s.updated = time.Now()

	return s.basicAuth
}

func (lr loginRequest) execute() (string, error) {
	jsonData, err := json.Marshal(lr)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%v/%v", lr.Host, sessionsPath)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := lr.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	var data map[string]string
	_ = json.Unmarshal(body, &data)

	return data["session_id"], nil
}

func createAuthHeader(sessionId string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:session", sessionId)))
}
