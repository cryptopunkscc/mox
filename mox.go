package mox

import (
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnd"
	"github.com/cryptopunkscc/go-xmppc/bot"
	"github.com/cryptopunkscc/go-xmppc/components/presence"
	"github.com/cryptopunkscc/mox/adminbot"
	"github.com/cryptopunkscc/mox/services"
	"github.com/cryptopunkscc/mox/xmpp"
	"github.com/cryptopunkscc/mox/xmpp/money"
	"log"
	"time"
)

type Mox struct {
	xmpp      *xmpp.XMPP
	ln        bitcoin.LightningClient
	bot       *bot.Bot
	lightning *services.LightningService

	quit chan bool
}

func New(cfg *Config) (*Mox, error) {
	lnc, _ := lnd.Connect(cfg.LND)
	mox := &Mox{
		quit: make(chan bool),
		xmpp: xmpp.NewXMPP(cfg.XMPP),
		ln:   lnc,
	}
	mox.lightning = &services.LightningService{lnc}
	mox.bot = bot.New(&adminbot.Engine{
		Lightning:  mox.lightning,
		Presence:   &mox.xmpp.Presence,
		Money:      mox.xmpp.Money,
		RosterComp: &mox.xmpp.Roster,
	})

	mox.xmpp.Money.InvoiceRequestHandler = func(req *money.InvoiceRequest) {
		log.Printf("%s asks for a %d SAT invoice!", req.JID, req.Amount.Sat())
		invoice := mox.lightning.IssueInvoice(req.Amount, "", time.Hour)
		req.SendInvoice(invoice.PaymentRequest)
	}

	mox.xmpp.Money.InvoiceHandler = func(inv *money.Invoice) {
		binvoice := mox.lightning.Decode(inv.Invoice)
		log.Printf("%s sent us an invoice for %d! Paying!", inv.JID, binvoice.Amount.Sat())
		err := mox.lightning.PayInvoice(inv.Invoice)
		if err != nil {
			log.Println("Error paying invoice:", err)
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
