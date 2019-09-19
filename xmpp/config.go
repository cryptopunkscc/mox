package xmpp

import "errors"

type Config struct {
	JID      string
	Password string
}

func (cfg *Config) Validate() error {
	if cfg.JID == "" {
		return errors.New("JID missing")
	}
	if cfg.Password == "" {
		return errors.New("Password missing")
	}
	return nil
}
