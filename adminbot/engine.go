package adminbot

import (
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-xmpp"
	chatbot "github.com/cryptopunkscc/go-xmppc/bot"
	"github.com/cryptopunkscc/go-xmppc/components/presence"
	"github.com/cryptopunkscc/mox/contacts"
	"github.com/cryptopunkscc/mox/payments"
	"time"
)

type Engine struct {
	Presence        *presence.Presence
	PaymentsService *payments.Service
	Contacts        *contacts.Service
}

func (e *Engine) Balance(ctx *chatbot.Context) {
	balance := e.PaymentsService.Balance()
	ctx.Reply("Your balance is %d SAT", balance.Sat())
}

func (e *Engine) Status(ctx *chatbot.Context, status string) {
	e.Presence.SetStatus(status)
	ctx.Reply("Status set.")
}

func (e *Engine) Issue(ctx *chatbot.Context, sats int, memo string) {
	if sats <= 0 {
		ctx.Reply("Amount must be greater than 0.")
		return
	}
	invoice := e.PaymentsService.IssueInvoice(bitcoin.Sat(int64(sats)), memo, 24*time.Hour)
	ctx.Reply("Here's your invoice:\n%s", invoice.PaymentRequest)
}

func (e *Engine) Pay(ctx *chatbot.Context, invoice string) {
	err := e.PaymentsService.PayInvoice(invoice)
	if err == nil {
		ctx.Reply("Invoice paid!")
		return
	}
	ctx.Reply("Failed to pay the invoice: %s", err.Error())
}

func (e *Engine) Send(ctx *chatbot.Context, jid string, amount int) {
	err := e.PaymentsService.SendBitcoin(xmpp.JID(jid), bitcoin.Sat(int64(amount)))

	if err != nil {
		ctx.Reply("Error sending money: %s", err.Error())
		return
	}

	ctx.Reply("Payment sent!")
}

func (e *Engine) List(ctx *chatbot.Context) {
	list := e.Contacts.AvailableContacts()

	for _, c := range list {
		var online, status string
		if c.Online {
			online = "*"
		}
		if c.Status != "" {
			status = "(" + c.Status + ")"
		}
		ctx.Reply("[%s%s] %s %s", c.JID, online, c.Name, status)
	}
}

func (e *Engine) Info(ctx *chatbot.Context) {
	me := e.Contacts.Me()
	ctx.Reply("%s (%s)", me.JID, me.Status)
}

func (e *Engine) Add(ctx *chatbot.Context, jid string, name string) {
	err := e.Contacts.AddContact(xmpp.JID(jid), name)
	if err != nil {
		ctx.Reply("Failed to add contact: %s", err)
		return
	}
	ctx.Reply("Contact added.")
}

func (e *Engine) Remove(ctx *chatbot.Context, jid string) {
	err := e.Contacts.RemoveContact(xmpp.JID(jid))
	if err != nil {
		ctx.Reply("Failed to remove contact: %s", err)
		return
	}
	ctx.Reply("Contact removed.")
}

func (e *Engine) Help(ctx *chatbot.Context, topic string) {
	switch topic {
	case "":
		ctx.Reply(
			"Commands:\n" +
				"status <status_text> - set XMPP status\n" +
				"balance - check your lightning balance\n" +
				"issue <amount> [memo] - issue a lightning invoice\n" +
				"pay <invoice> - pay a lightning invoice")
	case "status":
		ctx.Reply("Usage: status <status_text>")
	default:
		ctx.Reply("unknown help topic")
	}
}
