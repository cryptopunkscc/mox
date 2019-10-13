package xmpp

import (
	"errors"
	"github.com/cryptopunkscc/go-xmpp"
)

type Config struct {
	JID      xmpp.JID `yaml:"jid"`
	Password string   `yaml:"password"`
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
