package config

import (
	"os"
	"strconv"
)

type Config struct {
	SSHPort        int
	HostKeyPath    string
	DBPath         string
	Domain         string
	IdleTimeout    int
	MaxConnections int
	CaddyURL       string
}

func Load() *Config {
	return &Config{
		SSHPort:        getEnvInt("SSH_PORT", 2222),
		HostKeyPath:    getEnv("SSH_HOST_KEY_PATH", "ssh_host_key"),
		DBPath:         getEnv("DB_PATH", "poor-exe.db"),
		Domain:         getEnv("DOMAIN", "ssh.rnzlive.com"),
		IdleTimeout:    getEnvInt("IDLE_TIMEOUT", 1800), // 30 minutes
		MaxConnections: getEnvInt("MAX_CONNECTIONS", 100),
		CaddyURL:       getEnv("CADDY_URL", "http://localhost:2019"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
