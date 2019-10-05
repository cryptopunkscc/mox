package mox

import (
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnd"
	xmpp "github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmppc"
	"github.com/cryptopunkscc/go-xmppc/bot"
	"github.com/cryptopunkscc/go-xmppc/components/presence"
	"github.com/cryptopunkscc/mox/adminbot"
	"github.com/cryptopunkscc/mox/payments"
	"github.com/cryptopunkscc/mox/roster"
	xmppmox "github.com/cryptopunkscc/mox/xmpp"
	"log"
	"time"
)

type Mox struct {
	xmpp     *xmppmox.XMPP
	ln       bitcoin.LightningClient
	bot      *bot.Bot
	payments *payments.Service
	roster   *roster.Service

	quit chan bool
}

func (mox *Mox) Online(xmppc.Session) {
	mox.roster.Fetch()
}

func (mox *Mox) HandleStanza(xmpp.Stanza) {
}

func (mox *Mox) Offline(error) {
}

func New(cfg *Config) (*Mox, error) {
	lnc, _ := lnd.Connect(cfg.LND)
	mox := &Mox{
		quit: make(chan bool),
		xmpp: xmppmox.NewXMPP(cfg.XMPP),
		ln:   lnc,
	}
	mox.payments = &payments.Service{
		Component:       mox.xmpp.Payments,
		LightningClient: lnc,
	}
	mox.roster = &roster.Service{
		Roster:   mox.xmpp.Roster,
		Presence: mox.xmpp.Presence,
	}
	mox.bot = bot.New(&adminbot.Engine{
		PaymentsService: mox.payments,
		Presence:        mox.xmpp.Presence,
		Money:           mox.xmpp.Payments,
		RosterService:   mox.roster,
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

	mox.xmpp.Presence.UpdateHandler = func(update *presence.Update) {
		mox.roster.UpdatePresence(update)
	}

	mox.xmpp.Add(mox)
	mox.xmpp.Add(mox.bot)
	return mox, nil
}

func (mox *Mox) Run() {
	mox.xmpp.Connect()
	<-mox.quit
}
