package xmpp

import (
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmpp/ext/ping"
	"github.com/cryptopunkscc/go-xmpp/ext/presence"
	"github.com/cryptopunkscc/go-xmpp/ext/roster"
	"github.com/cryptopunkscc/mox/payments"
	"log"
	"time"
)

const reconnectInterval = 5 * time.Second
const defaultResource = "mox"

type XMPP struct {
	Ping     *ping.Ping
	Presence *presence.Presence
	Roster   *roster.Roster
	Payments *payments.Component

	cfg     *Config
	session xmpp.Session
	xmpp.Broadcast
}

func (x *XMPP) Online(s xmpp.Session) {
	x.session = s

	log.Println("XMPP bound to", s.JID())
}

func (x *XMPP) Offline(err error) {
	if err == nil {
		log.Println("Disconnected.")
	} else {
		log.Printf("Disconnected (error: %s). Reconnecting in %s...\n", err.Error(), reconnectInterval)
		go x.reconnect()
	}
}

func (x *XMPP) JID() xmpp.JID {
	if x.session == nil {
		return ""
	}
	return x.session.JID()
}

func (x *XMPP) HandleStanza(s xmpp.Stanza) {
}

func NewXMPP(cfg *Config) *XMPP {
	x := &XMPP{cfg: cfg}
	x.Add(x)

	// Setup ping component
	x.Ping = &ping.Ping{}
	x.Add(x.Ping)

	// Set up presence component
	x.Presence = &presence.Presence{}
	x.Add(x.Presence)

	// Set up roster component
	x.Roster = &roster.Roster{}
	x.Add(x.Roster)

	// Set up payments component
	x.Payments = &payments.Component{}
	x.Add(x.Payments)

	return x
}

func (x *XMPP) Connect() error {
	jid := x.cfg.JID
	if jid.Resource() == "" {
		jid = jid + "/" + defaultResource
	}
	xmppConfig := &xmpp.Config{
		JID:      jid,
		Password: x.cfg.Password,
	}

	return xmpp.Open(&x.Broadcast, xmppConfig)
}

func (x *XMPP) reconnect() {
	<-time.After(reconnectInterval)
	log.Println("Reconnecting...")
	if err := x.Connect(); err != nil {
		log.Printf("Error reconnecting: %s. Retrying in %s.\n", err.Error(), reconnectInterval)
		go x.reconnect()
	}
}
