package money

import (
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmppc"
)

var _ xmppc.Handler = &Money{}

type InvoiceRequestHandler func(*InvoiceRequest)
type InvoiceHandler func(*Invoice)

type InvoiceRequest struct {
	JID    xmpp.JID
	Amount bitcoin.Amount
	money  *Money
}

type Invoice struct {
	JID     xmpp.JID
	Invoice string
}

func (ir *InvoiceRequest) SendInvoice(invoice string) {
	ir.money.SendInvoice(ir.JID, invoice)
}

type Money struct {
	session xmppc.Session
	InvoiceRequestHandler
	InvoiceHandler
}

func (m *Money) RequestInvoice(jid xmpp.JID, amount bitcoin.Amount) {
	xm := &xmppMoney{
		Request: &xmppRequest{
			Amount: int(amount.Sat()),
		},
	}
	xmsg := &xmpp.Message{
		To:   jid,
		Type: "normal",
	}
	xmsg.AddChild(xm)
	m.session.Write(xmsg)
}

func (m *Money) SendInvoice(jid xmpp.JID, invoice string) {
	msg := &xmpp.Message{
		To:   jid,
		Type: "normal",
	}
	msg.AddChild(&xmppMoney{
		Invoice: &xmppInvoice{
			Encoded: invoice,
		},
	})
	m.session.Write(msg)
}

func (m *Money) handleMoney(msg *xmpp.Message) {
	money := msg.Child(&xmppMoney{}).(*xmppMoney)

	if money.Request != nil {
		if m.InvoiceRequestHandler != nil {
			m.InvoiceRequestHandler(&InvoiceRequest{
				JID:    msg.From,
				Amount: bitcoin.Sat(int64(money.Request.Amount)),
				money:  m,
			})
		}
	}
	if money.Invoice != nil {
		if m.InvoiceHandler != nil {
			m.InvoiceHandler(&Invoice{
				JID:     msg.From,
				Invoice: money.Invoice.Encoded,
			})
		}
	}
}

func (m *Money) HandleMessage(msg *xmpp.Message) {
	if _, ok := msg.Child(&xmppMoney{}).(*xmppMoney); ok {
		m.handleMoney(msg)
	}
}

func (m *Money) HandleStanza(s xmpp.Stanza) {
	xmppc.HandleStanza(m, s)
}

func (m *Money) Online(s xmppc.Session) {
	m.session = s
}

func (m *Money) Offline(error) {

}
