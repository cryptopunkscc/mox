package mox

import (
	"log"
	"os"
	"time"

	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmppc"
	"github.com/cryptopunkscc/go-xmppc/bot"
	"github.com/cryptopunkscc/go-xmppc/components/presence"
	"github.com/cryptopunkscc/mox/adminbot"
	"github.com/cryptopunkscc/mox/contacts"
	"github.com/cryptopunkscc/mox/payments"
	"github.com/cryptopunkscc/mox/wallet"
	xmppmox "github.com/cryptopunkscc/mox/xmpp"
)

type Mox struct {
	xmpp     *xmppmox.XMPP
	bot      *bot.Bot
	payments *payments.Service
	contacts *contacts.Service
	wallet   *wallet.Service

	quit chan bool
}

func (mox *Mox) Online(s xmppc.Session) {
	mox.contacts.SetJID(s.JID())
	mox.contacts.Fetch()
}

func (mox *Mox) HandleStanza(s xmpp.Stanza) {
}

func (mox *Mox) Offline(error) {
}

func New(cfg *Config) (*Mox, error) {
	var err error

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	xmppClient := xmppmox.NewXMPP(cfg.XMPP)

	mox := &Mox{
		quit: make(chan bool),
		xmpp: xmppClient,
	}

	// Set up wallet service
	mox.wallet, err = wallet.New(cfg.Wallet)
	if err != nil {
		return nil, err
	}

	// Set up payments service
	mox.payments = &payments.Service{
		Component: mox.xmpp.Payments,
		Wallet:    mox.wallet,
	}

	// Set up contacts service
	mox.contacts = &contacts.Service{
		Roster:   mox.xmpp.Roster,
		Presence: mox.xmpp.Presence,
	}

	// Set up admin bot
	mox.bot = bot.New(&adminbot.Engine{
		PaymentsService: mox.payments,
		Presence:        mox.xmpp.Presence,
		Contacts:        mox.contacts,
	})

	//mox.ln.SetInvoiceHandler(func(i *bitcoin.Invoice) {
	//	if i.Paid() {
	//		bytes, _ := json.MarshalIndent(i, "", "  ")
	//		fmt.Println(string(bytes))
	//		_ = mox.bot.Send("arashi@cryptopunks.cc", "Received %d SAT", i.Amount.Sat())
	//	}
	//})

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
		mox.contacts.UpdatePresence(update)
	}

	mox.xmpp.Add(mox)
	mox.xmpp.Add(mox.bot)
	return mox, nil
}

func (mox *Mox) Run() {
	_ = mox.xmpp.Connect()
	<-mox.quit
}
