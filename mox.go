package mox

import (
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnd"
	"github.com/cryptopunkscc/go-xmppc/bot"
	"github.com/cryptopunkscc/go-xmppc/components/presence"
	"github.com/cryptopunkscc/mox/adminbot"
	"github.com/cryptopunkscc/mox/payments"
	"github.com/cryptopunkscc/mox/xmpp"
	"log"
	"time"
)

type Mox struct {
	xmpp     *xmpp.XMPP
	ln       bitcoin.LightningClient
	bot      *bot.Bot
	payments *payments.Service

	quit chan bool
}

func New(cfg *Config) (*Mox, error) {
	lnc, _ := lnd.Connect(cfg.LND)
	mox := &Mox{
		quit: make(chan bool),
		xmpp: xmpp.NewXMPP(cfg.XMPP),
		ln:   lnc,
	}
	mox.payments = &payments.Service{
		Component:       mox.xmpp.Payments,
		LightningClient: lnc,
	}
	mox.bot = bot.New(&adminbot.Engine{
		Payments:   mox.payments,
		Presence:   &mox.xmpp.Presence,
		Money:      mox.xmpp.Payments,
		RosterComp: &mox.xmpp.Roster,
	})

	mox.xmpp.Payments.InvoiceRequestHandler = func(req *payments.InvoiceRequest) {
		invoice := mox.payments.IssueInvoice(req.Amount, "", time.Hour)
		err := req.SendInvoice(invoice.PaymentRequest)
		if err != nil {
			log.Println("Error sending invoice:", err)
		}
	}

	mox.xmpp.Presence.RequestHandler = func(request *presence.Request) {
		request.Allow()
	}

	mox.xmpp.Add(mox.bot)
	return mox, nil
}

func (mox *Mox) Run() {
	mox.xmpp.Connect()
	<-mox.quit
}
