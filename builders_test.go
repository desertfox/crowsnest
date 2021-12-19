package main

import (
	"testing"
)

//Bypass empty string checks
func newTestReqParams(u, p, c string) reqParams {
	return reqParams{
		Username:   u,
		Password:   p,
		ConfigPath: c,
	}
}

func Test_newReqParams(t *testing.T) {
	tests := []struct {
		name, username, password, config string
		want                             reqParams
		wantErr                          error
	}{
		{"missing username", "", "PASS", "$CONFIG", reqParams{}, errReqParams},
		{"missing password", "USER", "", "$CONFIG", reqParams{}, errReqParams},
		{"missing config", "USER", "PASS", "", reqParams{}, errReqParams},
		{"all params present", "USER", "PASS", "CONFIG", reqParams{"USER", "PASS", "CONFIG"}, nil},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newReqParams(tt.username, tt.password, tt.config)
			if got != tt.want {
				t.Fatalf("reqParams not expected, test: %v, got: %v", tt, got)
			}

			if gotErr != tt.wantErr {
				t.Fatalf("errors expected but not returned, test: %v", tt)
			}
		})
	}
}

/*
func Test_BuildConfig(t *testing.T) {
	emptyFile, err := ioutil.TempFile("", "test*")
	if err != nil {
		t.Fatalf("error constructing test cases %v", err)
	}

	tests := []struct {
		name      string
		reqParams reqParams
		want      config
		wantErr   error
	}{
		{"cannot find file", newTestReqParams("USER", "PASS", "TESTCONFIG"), config{}, errors.New("open TESTCONFIG: The system cannot find the file specified.")},
		{"empty config", newTestReqParams("USER", "PASS", emptyFile.Name()), config{}, nil},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := buildConfig(tt.reqParams)
			if got.Host != tt.want.Host {
				t.Fatalf("buildConfig not expected, got: %v, want: %v", got, tt.want)
			}

			if tt.wantErr != nil && gotErr.Error() == tt.wantErr.Error() {
				t.Fatalf("buildConfig errors expected but not returned, got: %v, want: %v", gotErr, tt.wantErr)
			}

		})
	}

}
*/
