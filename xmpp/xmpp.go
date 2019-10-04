package xmpp

import (
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmppc"
	"github.com/cryptopunkscc/go-xmppc/components/ping"
	"github.com/cryptopunkscc/go-xmppc/components/presence"
	"github.com/cryptopunkscc/go-xmppc/components/roster"
	"github.com/cryptopunkscc/mox/payments"
	"log"
	"time"
)

const reconnectInterval = 5 * time.Second

type XMPP struct {
	xmppc.Broadcast
	cfg      *Config
	session  xmppc.Session
	Presence presence.Presence
	Roster   roster.Roster
	Payments *payments.Component
}

func (xmpp *XMPP) Online(s xmppc.Session) {
	xmpp.session = s

	log.Println("Connected as", s.JID())
}

func (xmpp *XMPP) Offline(err error) {
	if err == nil {
		log.Println("Disconnected.")
	} else {
		log.Printf("Disconnected (error: %s). Reconnecting in %s...\n", err.Error(), reconnectInterval)
		go xmpp.reconnect()
	}
}

func (xmpp *XMPP) HandleStanza(s xmpp.Stanza) {
}

func NewXMPP(cfg *Config) *XMPP {
	xmpp := &XMPP{
		cfg:      cfg,
		Payments: &payments.Component{},
	}
	xmpp.Add(xmpp)
	xmpp.Add(&ping.Ping{})
	xmpp.Add(&xmpp.Presence)
	xmpp.Add(&xmpp.Roster)
	xmpp.Add(xmpp.Payments)

	return xmpp
}

func (xmpp *XMPP) Connect() error {
	xmppConfig := &xmppc.Config{
		JID:      xmpp.cfg.JID,
		Password: xmpp.cfg.Password,
	}

	return xmppc.Open(&xmpp.Broadcast, xmppConfig)
}

func (xmpp *XMPP) reconnect() {
	<-time.After(reconnectInterval)
	log.Println(("Reconnecting..."))
	if err := xmpp.Connect(); err != nil {
		log.Printf("Error reconnecting: %s. Retrying in %s.\n", err.Error(), reconnectInterval)
		go xmpp.reconnect()
	}
}
