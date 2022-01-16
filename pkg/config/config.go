package config

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
		Username:  checkEnvStr(loadEnv("username"), true),
		Password:  checkEnvStr(loadEnv("password"), true),
		Path:      checkEnvStr(loadEnv("path"), true),
		DelayJobs: checkEnvInt(loadEnv("delay"), false),
	}
}

func loadEnv(s string) string {
	return os.Getenv(buildEnvStr(s))
}

func checkEnvStr(s string, required bool) string {
	if s == "" && !required {
		return ""
	}

	if s == "" && required {
		log.Fatalf("missing enviromental variable: %s", buildEnvStr(s))
	}

	return s
}

func checkEnvInt(s string, required bool) int {
	if s == "" && !required {
		return 0
	}

	if s == "" && required {
		log.Fatalf("missing enviromental variable: %s", buildEnvStr(s))
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("cannot convert string to int: %s=%s", buildEnvStr(s), s)
	}

	return i
}

func buildEnvStr(s string) string {
	return prefix + strings.ToUpper(s)
}
