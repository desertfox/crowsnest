package main

import (
	"os"
	"testing"
)

func mockENV(c, u, p string) func() {
	cn := "CROWSNEST_"

	os.Setenv(cn+c, c)
	os.Setenv(cn+u, u)
	os.Setenv(cn+p, p)

	return func() {
		os.Unsetenv(cn + c)
		os.Unsetenv(cn + u)
		os.Unsetenv(cn + p)
	}
}

func TestEmptyEnv_BuildConfigFromENV(t *testing.T) {
	_, got := buildConfigFromENV()

	if got == nil {
		t.Error("No error returned")
	}

}

func TestEnv_BuildConfigFromENV(t *testing.T) {
	clearEnv := mockENV("CONFIG", "USERNAME", "PASSWORD")
	defer clearEnv()

	_, got := buildConfigFromENV()
	if got != nil {
		t.Errorf("Error present: %v", got)
	}
}
