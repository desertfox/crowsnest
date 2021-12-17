package main

import (
	"testing"
	"time"
)

func TestNewAuth(t *testing.T) {

	want := auth{"basicAuth", time.Now()}
	got := newAuth("basicAuth")

	if want.basicAuth != got.basicAuth {
		t.Errorf("basicAuth text does not match. got:%v want:%v", got.basicAuth, want.basicAuth)
	}

}
