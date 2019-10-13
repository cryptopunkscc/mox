package mox

import (
	"errors"
	"github.com/cryptopunkscc/mox/wallet"
	"io/ioutil"

	"github.com/cryptopunkscc/go-bitcoin/lnd"
	"github.com/cryptopunkscc/mox/rpcserver"
	"github.com/cryptopunkscc/mox/xmpp"
	"gopkg.in/yaml.v2"
)

type Config struct {
	RPC    *rpcserver.Config `yaml:"rpc"`
	XMPP   *xmpp.Config      `yaml:"xmpp"`
	LND    *lnd.Config       `yaml:"lnd"`
	Wallet *wallet.Config    `yaml:"wallet"`
}

func LoadConfig(configFile string) *Config {
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	cfg := &Config{}
	err = yaml.Unmarshal(bytes, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func (cfg *Config) Validate() error {
	// Validate XMPP config
	if err := cfg.XMPP.Validate(); err != nil {
		return err
	}

	// Validate RPC config
	if cfg.RPC != nil {
		if err := cfg.RPC.Validate(); err != nil {
			return err
		}
	}

	// Validate Wallet config
	if cfg.Wallet == nil {
		return errors.New("wallet config missing")
	}
	if err := cfg.Wallet.Validate(); err != nil {
		return err
	}

	return nil
}
