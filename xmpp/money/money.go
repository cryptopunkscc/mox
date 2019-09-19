package money

import (
	"fmt"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-rx"
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmppc"
)

type InvoiceRequest struct {
	From   string
	Amount int
	ID     string
}

type InvoiceRequestHandler func(*InvoiceRequest)
type InvoiceHandler func(string)

type Money struct {
	client               *xmppc.Client
	InvoiceRequestStream rx.SyncStream
	InvoiceStream        rx.SyncStream
}

func New(c *xmppc.Client) *Money {
	money := &Money{client: c}
	c.IQStream.Subscribe(money.handleIQ, nil, nil)
	return money
}

func (money *Money) RequestInvoice(to string, amount int, output rx.Stream) {
	iq := xmppc.MakeIQ(to, "get", "")

	iq.AddChild(&xmppInvoice{Amount: amount})

	money.client.Write(iq, rx.Pipe(func(res *xmpp.IQ) {
		if i, ok := res.Child("invoice").(*xmppInvoice); ok {
			fmt.Println("[XMPP-MONEY]", i.Data)
			money.InvoiceStream.Next(i)
			if output != nil {
				output.Next(i.Data)
			}
		}

	}))
}

func (money *Money) SendInvoice(to string, id string, invoice *bitcoin.Invoice) {
	r := xmppc.MakeIQ(to, "result", "")
	r.ID = id
	r.AddChild(&xmppInvoice{
		Data: invoice.PaymentRequest,
	})
	money.client.Write(r, nil)
}

func (money *Money) handleInvoiceRequest(s *xmpp.IQ, i *xmppInvoice) {
	req := &InvoiceRequest{
		ID:     s.ID,
		From:   s.From,
		Amount: i.Amount,
	}

	money.InvoiceRequestStream.Next(req)
}

func (money *Money) handleIQ(s *xmpp.IQ) {
	if s.Type != "get" {
		return
	}

	if i, ok := s.Child("invoice").(*xmppInvoice); ok {
		money.handleInvoiceRequest(s, i)
	}
}

func init() {
	xmpp.IQContext.Add(&xmppInvoice{})
}
