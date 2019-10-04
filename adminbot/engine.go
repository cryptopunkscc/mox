package adminbot

import (
	"fmt"
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-xmpp"
	chatbot "github.com/cryptopunkscc/go-xmppc/bot"
	"github.com/cryptopunkscc/go-xmppc/components/presence"
	"github.com/cryptopunkscc/go-xmppc/components/roster"
	"github.com/cryptopunkscc/mox/payments"
	"time"
)

type Engine struct {
	Presence   *presence.Presence
	RosterComp *roster.Roster
	Money      *payments.Component
	Payments   *payments.Service
}

func (e *Engine) Balance(ctx *chatbot.Context) {
	balance := e.Payments.Balance()
	ctx.Reply("Your balance is %d SAT", balance.Sat())
}

func (e *Engine) Status(ctx *chatbot.Context, status string) {
	if status == "" {
		e.Help(ctx, "status")
		return
	}
	e.Presence.SetStatus(status)
}

func (e *Engine) Issue(ctx *chatbot.Context, sats int, memo string) {
	if sats <= 0 {
		ctx.Reply("Amount must be greater than 0.")
		return
	}
	invoice := e.Payments.IssueInvoice(bitcoin.Sat(int64(sats)), memo, 24*time.Hour)
	ctx.Reply("Here's your invoice:\n%s", invoice.PaymentRequest)
}

func (e *Engine) Pay(ctx *chatbot.Context, invoice string) {
	err := e.Payments.PayInvoice(invoice)
	if err == nil {
		ctx.Reply("Invoice paid!")
		return
	}
	ctx.Reply("Failed to pay the invoice: %s", err.Error())
}

func (e *Engine) Send(ctx *chatbot.Context, jid string, amount int) {
	err := e.Payments.SendBitcoin(xmpp.JID(jid), bitcoin.Sat(int64(amount)))

	if err != nil {
		ctx.Reply("Error sending money: %s", err.Error())
		return
	}

	ctx.Reply("Payment sent!")
}

func (e *Engine) Roster(ctx *chatbot.Context) {
	e.RosterComp.FetchRoster(func(items []*roster.RosterItem) {
		res := "Roster:\n"
		for _, i := range items {
			res = res + fmt.Sprintf("%s <%s> [%s]\n", i.Name, i.JID, i.Subscription)
		}
		ctx.Reply(res)
	})
}

func (e *Engine) AddRoster(ctx *chatbot.Context, jid string, name string) {
	e.RosterComp.Add(xmpp.JID(jid), name)
	ctx.Reply("Added to roster.")
}

func (e *Engine) DelRoster(ctx *chatbot.Context, jid string) {
	e.RosterComp.Remove(xmpp.JID(jid))
	ctx.Reply("Removed from roster.")
}

func (e *Engine) Subscribe(ctx *chatbot.Context, jid string) {
	e.Presence.Subscribe(xmpp.JID(jid))
	ctx.Reply("Subscription request sent!")
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
