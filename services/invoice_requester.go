package services

import (
	"github.com/cryptopunkscc/go-rx"
	"github.com/cryptopunkscc/mox/xmpp"
)

type InvoiceRequster struct {
	xmpp *xmpp.XMPP
}

func NewInvoiceRequester(j *xmpp.XMPP) *InvoiceRequster {
	return &InvoiceRequster{
		xmpp: j,
	}
}

func (r *InvoiceRequster) Request(jid string, amount int, output rx.Stream) {
	r.xmpp.Money.RequestInvoice(jid, amount, output)
}
