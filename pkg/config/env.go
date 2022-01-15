package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

var prefix string = "CROWSNEST_"

type Env struct {
	Username   string
	Password   string
	ConfigPath string
	DelayJobs  int
}

func (e *Env) GetEnv() {
	e.Username = loadStrFromEnv("username", true)
	e.Password = loadStrFromEnv("password", true)
	e.ConfigPath = loadStrFromEnv("config", true)
	e.DelayJobs = loadIntFromEnv("delay", false)
}

func loadStrFromEnv(key string, required bool) string {
	value := os.Getenv(prefix + strings.ToUpper(key))

	if value == "" && !required {
		return ""
	}

	if value == "" && required {
		log.Fatalf("missing enviromental variable: %s", prefix+strings.ToUpper(key))
	}

	return value
}

func loadIntFromEnv(key string, required bool) int {
	value := os.Getenv(prefix + strings.ToUpper(key))

	if value == "" && !required {
		return 0
	}

	if value == "" && required {
		log.Fatalf("missing enviromental variable: %s", prefix+strings.ToUpper(key))
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("cannot convert string to int: %s=%s", prefix+strings.ToUpper(key), value)
	}

	return i
}
