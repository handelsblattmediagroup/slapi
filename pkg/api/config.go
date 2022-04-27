package api

import (
	"fmt"
	"os"
	"strconv"
)

var DefaultConfig Config

func GetDefaultConfig() *Config {
	return &Config{
		ListenAddr: fmt.Sprintf(":%s", os.Getenv("PORT")),
		Github: GithubConfig{
			PrivateKey:    os.Getenv("GITHUB_PRIVATE_KEY"),
			ClientID:      os.Getenv("GITHUB_CLIENT_ID"),
			IntegrationID: getenvInt("GITHUB_APP_ID"),
		},
	}
}

type Config struct {
	Github     GithubConfig
	ListenAddr string
}

type GithubConfig struct {
	PrivateKey    string
	ClientID      string
	IntegrationID int64
}

func getenvInt(name string) int64 {
	valueStr := os.Getenv(name)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0
	}
	return value
}
