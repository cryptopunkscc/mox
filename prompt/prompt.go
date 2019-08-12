package prompt

import (
	"fmt"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-rx"
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmppc/components/disco"
	"github.com/cryptopunkscc/mox/jabber"
	"github.com/cryptopunkscc/mox/services"
)

type Prompt struct {
	xmpp        *jabber.Jabber
	ln          bitcoin.LightningClient
	moneySender *services.MoneySender
}

func New(j *jabber.Jabber, ln bitcoin.LightningClient, ms *services.MoneySender) *Prompt {
	p := &Prompt{
		xmpp:        j,
		ln:          ln,
		moneySender: ms,
	}
	go p.CommandLine()
	return p
}

func (prompt *Prompt) CommandLine() {
	var cmd string
	for {
		fmt.Scanf("%s", &cmd)
		switch cmd {
		case "hello":
			fmt.Println("world!")
		case "add":
			var jid string
			fmt.Scanf("%s", &jid)
			prompt.xmpp.Roster.Subscribe(xmpp.JID(jid))
		case "remove":
			var jid string
			fmt.Scanf("%s", &jid)
			prompt.xmpp.Roster.Unsubscribe(xmpp.JID(jid))
		case "msg":
			var jid, msg string
			fmt.Scanf("%s %s", &jid, &msg)
			prompt.xmpp.Chat.SendMessage(jid, msg)
		case "status":
			var status string
			fmt.Scanf("%s", &status)
			prompt.xmpp.Presence.SetStatus(status)
		case "balance":
			fmt.Println("Your balance is", prompt.ln.Balance())
		case "send":
			var jid string
			var amt int
			fmt.Scanf("%s %d", &jid, &amt)
			prompt.moneySender.Send(jid, amt)
		case "query":
			var jid string
			fmt.Scanf("%s", &jid)
			prompt.xmpp.Disco.Query(jid, rx.SyncPipe(printDiscoInfo))
		}
	}
}

func printDiscoInfo(info *disco.Info) {
	fmt.Println("Features:")
	for _, f := range info.Features {
		fmt.Println("-", f)
	}
	fmt.Println("Identities:")
	for _, i := range info.Identities {
		fmt.Println("-", i.Name, i.Category, i.Type)
	}
}
