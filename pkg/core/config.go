package core

import "os"

func GetConfigDefaults() *Config {
	return &Config{
		ListenAddr:    getenvStringDefault("SLAPI_LISTEN_ADDR", ":8080"),
		AllowdOrigins: os.Getenv("SLAPI_ALLOWED_ORIGINS"),
		LogLevel:      getenvStringDefault("SLAPI_LOG_LEVEL", "info"),
	}
}

type Config struct {
	ListenAddr    string
	LogLevel      string
	AllowdOrigins string
}

func getenvStringDefault(name string, def string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	return value
}
