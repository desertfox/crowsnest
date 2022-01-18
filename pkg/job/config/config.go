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

func LoadFromEnv() *Config {
	return &Config{
		Username:  checkEnvStr("username", true),
		Password:  checkEnvStr("password", true),
		Path:      checkEnvStr("config", true),
		DelayJobs: checkEnvInt("delay", false),
	}
}

func loadEnv(s string) string {
	return os.Getenv(buildEnvStr(s))
}

func buildEnvStr(s string) string {
	return prefix + strings.ToUpper(s)
}

func checkEnvStr(s string, required bool) string {
	str := loadEnv(s)

	if str == "" && !required {
		return ""
	}

	if str == "" && required {
		log.Fatalf("missing enviromental variable: %s", buildEnvStr(s))
	}

	return str
}

func checkEnvInt(s string, required bool) int {
	str := loadEnv(s)

	if str == "" && !required {
		return 0
	}

	if str == "" && required {
		log.Fatalf("missing enviromental variable: %s", buildEnvStr(s))
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("cannot convert string to int: %s=%s", buildEnvStr(s), str)
	}

	return i
}
