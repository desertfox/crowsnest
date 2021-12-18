package main

import (
	"errors"
	"testing"
)

func newTestReqParams(u, p, c string) reqParams {
	return reqParams{
		Username:   u,
		Password:   p,
		ConfigPath: c,
	}
}

func mockConfig() {

}

func Test_newReqParams(t *testing.T) {
	tests := []struct {
		params []string
		pass   bool
	}{
		{[]string{"", "PASS", "$CONFIG"}, false},
		{[]string{"USER", "", "$CONFIG"}, false},
		{[]string{"USER", "PASS", ""}, false},
		{[]string{"USER", "PASS", "CONFIG"}, true},
	}

	for _, tt := range tests {
		_, got := newReqParams(tt.params[0], tt.params[1], tt.params[2])
		if !tt.pass && got == nil {
			t.Fatalf("expected error but was nil, test: %v, got: %v", tt, got)
			break
		}

		if tt.pass && got != nil {
			t.Fatalf("expected pass but was nil, test: %v, got: %v", tt, got)
			break
		}
	}

}

func TestEmptyConfig_BuildConfigFromENV(t *testing.T) {
	rp := newTestReqParams("USER", "PASS", "TESTCONFIG")
	_, got := buildConfigFromENV(rp)

	want := errors.New("open TESTCONFIG: The system cannot find the file specified.")

	if got == nil {
		t.Fatal("Expected error")
	}

	if got.Error() != want.Error() {
		t.Fatalf("Error does not match expect, got: %v, want: %v", got, want)
	}

}

func TestEnv_BuildConfigFromENV(t *testing.T) {
	rp := newTestReqParams("USER", "PASS", "TESTCONFIG")
	_, got := buildConfigFromENV(rp)

	if got != nil {
		t.Errorf("Error present: %v", got)
	}
}
