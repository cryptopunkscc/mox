package wallet

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

const defaultBackend = "backend"
const defaultLNDHost = "localhost"
const defaultLNDPort = 10009

type Config struct {
	Backend  string `yaml:"backend"`
	LNDDir   string `yaml:"lnddir"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Macaroon string `yaml:"macaroon"`
	Cert     string `yaml:"cert"`
}

func (cfg *Config) Validate() error {
	if cfg.getBackend() != "backend" {
		return errors.New("wallet: only backend backend is supported")
	}
	if cfg.getHost() == "" {
		return errors.New("wallet: host missing")
	}
	if (cfg.getPort() < 1) || (cfg.Port > 65535) {
		return errors.New("wallet: invalid port")
	}
	if _, err := os.Stat(cfg.getMacaroonPath()); os.IsNotExist(err) {
		return errors.New("wallet: macaroon file does not exist")
	}
	if _, err := os.Stat(cfg.getCertPath()); os.IsNotExist(err) {
		return errors.New("wallet: cert file does not exist")
	}
	return nil
}

func (cfg *Config) getBackend() string {
	if cfg.Backend != "" {
		return cfg.Backend
	}
	return defaultBackend
}

func (cfg *Config) getHost() string {
	if cfg.Host != "" {
		return cfg.Host
	}
	return defaultLNDHost
}

func (cfg *Config) getPort() int {
	if cfg.Port != 0 {
		return cfg.Port
	}
	return defaultLNDPort
}

func (cfg *Config) getLNDDir() string {
	if cfg.LNDDir != "" {
		return cfg.LNDDir
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".backend")
	}
	return ""
}

func (cfg *Config) getMacaroonPath() string {
	if cfg.Macaroon != "" {
		return cfg.Macaroon
	}
	if lnddir := cfg.getLNDDir(); lnddir != "" {
		return filepath.Join(lnddir, "data", "chain", "bitcoin", "testnet", "admin.macaroon")
	}
	return ""
}

func (cfg *Config) getCertPath() string {
	if cfg.Cert != "" {
		return cfg.Cert
	}
	if lnddir := cfg.getLNDDir(); lnddir != "" {
		return filepath.Join(lnddir, "tls.cert")
	}
	return ""
}

func (cfg *Config) getMacaroon() []byte {
	path := cfg.getMacaroonPath()
	if path == "" {
		return nil
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return bytes
}

func (cfg *Config) getCert() []byte {
	path := cfg.getCertPath()
	if path == "" {
		return nil
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return bytes
}
