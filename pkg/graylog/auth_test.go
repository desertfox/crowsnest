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
		wantErr                        error
	}{
		{"missing username", "", "USER", "PASS", loginRequest{}, errMissingParam},
		{"missing password", "HOST", "", "PASS", loginRequest{}, errMissingParam},
		{"missing config", "HOST", "USER", "", loginRequest{}, errMissingParam},
		{"all params present", "HOST", "USER", "PASS", loginRequest{"HOST", "USER", "PASS", client}, nil},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newLoginRequest(tt.host, tt.username, tt.password, client)
			if *got != tt.want {
				t.Fatalf("got: %v, want: %v", got, tt.want)
			}

			if gotErr != tt.wantErr {
				t.Fatalf("got: %v, want: %v", gotErr, tt.wantErr)
			}
		})
	}
}
