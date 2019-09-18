package mox

import (
	"github.com/cryptopunkscc/go-bitcoin/lnd"
	"github.com/cryptopunkscc/go-dep"
	"github.com/cryptopunkscc/mox/chatbot"
	"github.com/cryptopunkscc/mox/jabber"
	"github.com/cryptopunkscc/mox/prompt"
	"github.com/cryptopunkscc/mox/rpcserver"
	"github.com/cryptopunkscc/mox/services"
)

type Mox struct {
	*dep.Context
	jabber *jabber.Jabber

	quit chan bool
}

func New(cfg *Config) *Mox {
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	mox := &Mox{
		Context: dep.NewContext(),
		quit:    make(chan bool),
	}

	mox.Add(cfg.Jabber)
	mox.Add(cfg.LND)
	mox.Add(cfg.RPC)

	mox.Make(jabber.NewJabber)
	mox.Make(lnd.Connect)
	mox.Make(services.NewInvoiceSender)
	mox.Make(services.NewInvoiceRequester)
	mox.Make(services.NewMoneySender)
	mox.Make(services.NewBalanceChecker)
	mox.Make(services.NewInvoiceDecoder)
	mox.Make(rpcserver.New, cfg.RPC)
	mox.Make(chatbot.New, cfg.Chatbot)
	mox.Make(prompt.New)

	return mox
}

func (mox *Mox) Run() {
	mox.Call(func(rpc *rpcserver.RPCServer, j *jabber.Jabber) {
		//go rpc.Run()
		if err := j.Connect(); err != nil {
			panic(err)
		}
	})

	<-mox.quit
}
