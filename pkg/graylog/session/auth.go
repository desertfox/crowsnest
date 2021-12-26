package session

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const sessionsPath string = "api/system/sessions"

var errMissingParam = errors.New("missing param")

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

func New(h, u, p string, httpClient *http.Client) (*session, error) {
	lr, err := newLoginRequest(h, u, p, httpClient)
	if err != nil {
		return &session{}, err
	}
	return &session{"", time.Now(), lr}, nil
}

func newLoginRequest(h, u, p string, httpClient *http.Client) (*loginRequest, error) {
	for _, s := range []string{h, u, p} {
		if s == "" {
			return &loginRequest{}, errMissingParam

		}
	}

	return &loginRequest{h, u, p, httpClient}, nil
}

func (s session) GetHost() string {
	return s.loginRequest.Host
}

func (s *session) GetHeader() string {
	//check if token is old
	if 1 == 1 {
		sessionId, err := s.loginRequest.execute()
		if err != nil {
			panic(err.Error())
		}

		s.basicAuth = createAuthHeader(sessionId)
		s.updated = time.Now()
	}
	// Token is good
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
