package jobs

import (
	"log"
	"os"
	"strconv"
	"strings"
)

var prefix string = "CROWSNEST_"

type Config struct {
	Username  string
	Password  string
	Path      string
	DelayJobs int
}

func LoadConfigFromEnv() *Config {
	return &Config{
		Username:  loadStrFromEnv("username", true),
		Password:  loadStrFromEnv("password", true),
		Path:      loadStrFromEnv("config", true),
		DelayJobs: loadIntFromEnv("delay", false),
	}
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
