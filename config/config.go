package config

import (
	"log"
	"os"
	"strings"
)

var prefix string = "CROWSNEST_"

type Config struct {
	Username string //Used by job
	Password string //""
	Path     string //Used by joblist
}

func Load() *Config {
	return &Config{
		Username: checkEnvStr("username", false),
		Password: checkEnvStr("password", false),
		Path:     checkEnvStr("config", true),
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
