package adminbot

import (
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-xmpp"
	chatbot "github.com/cryptopunkscc/go-xmpp/bot"
	"github.com/cryptopunkscc/go-xmpp/ext/presence"
	"github.com/cryptopunkscc/mox/contacts"
	"github.com/cryptopunkscc/mox/payments"
	"github.com/cryptopunkscc/mox/wallet"
	"time"
)

var _ chatbot.Engine = &Engine{}

type Engine struct {
	chatbot.ChatWriter
	Presence        *presence.Presence
	PaymentsService *payments.Service
	Contacts        *contacts.Service
	Wallet          *wallet.Service
}

func (e *Engine) Online(writer chatbot.ChatWriter) {
	e.ChatWriter = writer
}

func (e *Engine) Offline(error) {
	e.ChatWriter = nil
}

func (e *Engine) Balance(ctx *chatbot.Context) error {
	balance := e.PaymentsService.Balance()
	return ctx.Reply("Your balance is %d SAT", balance.Sat())
}

func (e *Engine) Status(ctx *chatbot.Context, status string) error {
	e.Presence.SetStatus(status)
	return ctx.Reply("Status set.")
}

func (e *Engine) Issue(ctx *chatbot.Context, sats int, memo string) error {
	if sats <= 0 {
		return ctx.Reply("Amount must be greater than 0.")
	}
	invoice := e.PaymentsService.IssueInvoice(bitcoin.Sat(int64(sats)), memo, 24*time.Hour)
	return ctx.Reply("Here's your invoice:\n%s", invoice.PaymentRequest)
}

func (e *Engine) Pay(ctx *chatbot.Context, invoice string) error {
	err := e.PaymentsService.PayInvoice(invoice)
	if err == nil {
		return ctx.Reply("Invoice paid!")
	}
	return ctx.Reply("Failed to pay the invoice: %s", err.Error())
}

func (e *Engine) Send(ctx *chatbot.Context, jid string, amount int) error {
	err := e.PaymentsService.SendBitcoin(xmpp.JID(jid), bitcoin.Sat(int64(amount)))

	if err != nil {
		return ctx.Reply("Error sending money: %s", err.Error())
	}

	return ctx.Reply("Payment sent!")
}

func (e *Engine) List(ctx *chatbot.Context) error {
	list := e.Contacts.Contacts(contacts.All)

	for _, c := range list {
		var online, status string
		if c.Online {
			online = "*"
		}
		if c.Status != "" {
			status = "(" + c.Status + ")"
		}
		_ = ctx.Reply("[%s%s] %s %s", c.JID, online, c.Name, status)
	}
	return nil
}

func (e *Engine) Info(ctx *chatbot.Context) error {
	me := e.Contacts.Me()
	return ctx.Reply("%s (%s)", me.JID, me.Status)
}

func (e *Engine) Add(ctx *chatbot.Context, jid string, name string) error {
	err := e.Contacts.AddContact(xmpp.JID(jid), name)
	if err != nil {
		return ctx.Reply("Failed to add contact: %s", err)
	}
	return ctx.Reply("Contact added.")
}

func (e *Engine) Remove(ctx *chatbot.Context, jid string) error {
	err := e.Contacts.RemoveContact(xmpp.JID(jid))
	if err != nil {
		return ctx.Reply("Failed to remove contact: %s", err)
	}
	return ctx.Reply("Contact removed.")
}

func (e *Engine) ChainBalance(ctx *chatbot.Context) error {
	balance := e.Wallet.ChainBalance()
	return ctx.Reply("Your on-chain balance is %d SAT", balance.Sat())
}

func (e *Engine) ChainAddress(ctx *chatbot.Context) error {
	addr, err := e.Wallet.NewAddress()
	if err != nil {
		return ctx.Reply("Couldn't get new address: %s", err)
	}
	return ctx.Reply(addr)
}

func (e *Engine) ChainSend(ctx *chatbot.Context, addr string, amount int, feeRate int) error {
	txid, err := e.Wallet.ChainSend(addr, bitcoin.Sat(int64(amount)), feeRate)
	if err != nil {
		return ctx.Reply("Failed to send money: %s", err)
	}
	return ctx.Reply("Money sent! Transaction id: %s", txid)
}

func (e *Engine) Help(ctx *chatbot.Context, topic string) error {
	switch topic {
	case "":
		return ctx.Reply(
			"Commands:\n" +
				"status <status_text> - set XMPP status\n" +
				"balance - check your lightning balance\n" +
				"issue <amount> [memo] - issue a lightning invoice\n" +
				"pay <invoice> - pay a lightning invoice")
	case "status":
		return ctx.Reply("Usage: status <status_text>")
	default:
		return ctx.Reply("unknown help topic")
	}
}
