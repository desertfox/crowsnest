package graylog

import (
	"net/http"
	"testing"
)

func Test_newLoginRequest(t *testing.T) {
	var client *http.Client = &http.Client{}

	tests := []struct {
		name, host, username, password string
		want                           loginRequest
	}{
		{"all params present", "HOST", "USER", "PASS", loginRequest{"HOST", "USER", "PASS", client}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := loginRequest{tt.host, tt.username, tt.password, client}
			if got != tt.want {
				t.Fatalf("got: %v, want: %v", got, tt.want)
			}

		})
	}
}
