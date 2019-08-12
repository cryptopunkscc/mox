package rpcserver

import (
	"errors"
	"fmt"
)

const (
	defaultBind = "127.0.0.1"
	defaultPort = 50000
)

type Config struct {
	Port int    `json:"port"`
	Bind string `json:"bind"`
}

func (cfg *Config) Validate() error {
	if (cfg.Port < 1) || (cfg.Port > 65535) {
		return errors.New("Invalid RPC port")
	}
	if cfg.Bind == "" {
		return errors.New("RPC bind address missing")
	}
	return nil
}

func (cfg *Config) bind() string {
	if cfg.Bind != "" {
		return cfg.Bind
	}
	return "127.0.0.1"
}

func (cfg *Config) port() int {
	if cfg.Port != 0 {
		return cfg.Port
	}
	return 50000
}

func (cfg *Config) bindAddress() string {
	return fmt.Sprintf("%s:%d", cfg.bind(), cfg.port())
}
