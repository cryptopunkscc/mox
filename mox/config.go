package mox

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/cryptopunkscc/mox/chatbot"
	"github.com/cryptopunkscc/mox/rpcserver"

	"github.com/cryptopunkscc/go-bitcoin/lnd"
	"github.com/cryptopunkscc/mox/jabber"
)

type Config struct {
	RPC     *rpcserver.Config `json:"rpc"`
	Jabber  *jabber.Config    `json:"jabber"`
	Chatbot *chatbot.Config   `json:"chatbot"`
	LND     *lnd.Config       `json:"lnd"`
}

func LoadConfig(configFile string) *Config {
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	cfg := &Config{}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func (cfg *Config) Validate() error {
	if err := cfg.Jabber.Validate(); err != nil {
		return err
	}
	if cfg.RPC != nil {
		if err := cfg.RPC.Validate(); err != nil {
			return err
		}
	}
	if cfg.LND == nil {
		return errors.New("LND config missing")
	}
	if err := cfg.LND.Validate(); err != nil {
		return err
	}
	return nil
}
