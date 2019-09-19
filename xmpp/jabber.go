package xmpp

import (
	"fmt"
	"log"
	"time"

	"github.com/cryptopunkscc/mox/xmpp/money"

	"github.com/cryptopunkscc/go-dep"
	"github.com/cryptopunkscc/go-xmppc"
	"github.com/cryptopunkscc/go-xmppc/components/caps"
	"github.com/cryptopunkscc/go-xmppc/components/chat"
	"github.com/cryptopunkscc/go-xmppc/components/disco"
	"github.com/cryptopunkscc/go-xmppc/components/ping"
	"github.com/cryptopunkscc/go-xmppc/components/presence"
	"github.com/cryptopunkscc/go-xmppc/components/roster"
)

type XMPP struct {
	*xmppc.Client
	*dep.Context
	Ping     *ping.Ping
	Chat     *chat.Chat
	Roster   *roster.Roster
	Money    *money.Money
	Presence *presence.Presence
	Disco    *disco.Disco
	Caps     *caps.Caps
}

func NewXMPP(cfg *Config) (*XMPP, error) {
	c := dep.NewContext()
	j := &XMPP{Context: c}

	c.Add(xmppc.NewClient(cfg.JID, cfg.Password))
	c.Make(ping.New)
	c.Make(chat.New)
	c.Make(roster.New)
	c.Make(money.New)
	c.Make(presence.New)
	c.Make(disco.New)
	c.Make(caps.New)
	c.Inject(j)

	j.Chat.MessageStream.Subscribe(func(msg *chat.Message) {
		log.Printf("[%s] %s\n", msg.From.Bare(), msg.Body)
	}, nil, nil)

	j.Presence.UpdateStream.Subscribe(func(u *presence.Update) {
		a := "offline"
		if u.Online {
			a = "online"
		}
		fmt.Printf("%s has gone %s (%s)\n", u.JID, a, u.Status)
	}, nil, nil)

	j.Presence.RequestStream.Subscribe(func(req *presence.Request) {
		log.Println("Auto-accepting subscription request from", req.JID)
		req.Allow()
		j.Presence.Subscribe(req.JID)
	}, nil, nil)

	j.Roster.RosterStream.Subscribe(func(items []*roster.RosterItem) {
		for _, i := range items {
			log.Println(i.JID, i.Name, i.Subscription)
		}
	}, nil, nil)

	// j.Roster.UpdateStream.Subscribe(func(update *roster.Update) {
	// 	log.Println(update.JID, "changed status", update.Priority, update.Status, update.Show)
	// }, nil, nil)

	return j, nil
}

func logLatency(l time.Duration) {
	log.Println("Ping", l)
}
