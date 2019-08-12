package services

import (
	"github.com/cryptopunkscc/go-rx"
	"github.com/cryptopunkscc/mox/jabber"
)

type InvoiceRequster struct {
	jabber *jabber.Jabber
}

func NewInvoiceRequester(j *jabber.Jabber) *InvoiceRequster {
	return &InvoiceRequster{
		jabber: j,
	}
}

func (r *InvoiceRequster) Request(jid string, amount int, output rx.Stream) {
	r.jabber.Money.RequestInvoice(jid, amount, output)
}
