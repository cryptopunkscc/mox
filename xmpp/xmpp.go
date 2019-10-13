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
const defaultResource = "mox"

type XMPP struct {
	Ping     *ping.Ping
	Presence *presence.Presence
	Roster   *roster.Roster
	Payments *payments.Component

	cfg     *Config
	session xmppc.Session
	xmppc.Broadcast
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

func (xmpp *XMPP) Connect() error {
	jid := xmpp.cfg.JID
	if jid.Resource() == "" {
		jid = jid + "/" + defaultResource
	}
	xmppConfig := &xmppc.Config{
		JID:      jid,
		Password: xmpp.cfg.Password,
	}

	return xmppc.Open(&xmpp.Broadcast, xmppConfig)
}

func (xmpp *XMPP) reconnect() {
	<-time.After(reconnectInterval)
	log.Println("Reconnecting...")
	if err := xmpp.Connect(); err != nil {
		log.Printf("Error reconnecting: %s. Retrying in %s.\n", err.Error(), reconnectInterval)
		go xmpp.reconnect()
	}
}
