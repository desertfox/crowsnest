package session

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

type LoginRequest struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	session  auth
}

type auth struct {
	basicAuth string
	updated   time.Time
}

func NewLoginRequest(h, u, p string) *LoginRequest {
	for i, s := range []string{h, u, p} {
		if s == "" {
			switch i {
			case 0:
				panic("Missing host variable")
			case 1:
				panic("Missing username variable")
			case 2:
				panic("Missing password variable")
			}
		}
	}
	return &LoginRequest{h, u, p, auth{}}
}

func (lr LoginRequest) GetHost() string {
	return lr.Host
}

func (lr LoginRequest) GetSessionHeader(httpClient *http.Client) string {
	//check if token is old
	if 1 == 0 {
		sessionId, err := lr.sessionIdRequest(httpClient)
		if err != nil {
			panic(err.Error())
		}

		lr.session.basicAuth = createAuthHeader(sessionId)
		lr.session.updated = time.Now()
	}
	// Token is good
	return lr.session.basicAuth
}

func (lr LoginRequest) sessionIdRequest(httpClient *http.Client) (string, error) {
	jsonData, err := json.Marshal(lr)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%v/%v", lr.Host, sessionsPath)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := httpClient.Do(request)
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
