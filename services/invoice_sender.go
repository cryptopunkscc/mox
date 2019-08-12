package services

import (
	"time"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/mox/jabber"
	"github.com/cryptopunkscc/mox/jabber/money"
)

type InvoiceSender struct {
	jabber *jabber.Jabber
	ln     bitcoin.LightningClient
}

func NewInvoiceSender(j *jabber.Jabber, ln bitcoin.LightningClient) *InvoiceSender {
	s := &InvoiceSender{
		jabber: j,
		ln:     ln,
	}
	s.jabber.Money.InvoiceRequestStream.Subscribe(s.onRequest, nil, nil)
	return s
}

func (s *InvoiceSender) onRequest(req *money.InvoiceRequest) {
	i := s.ln.CreateInvoice(bitcoin.Sat(int64(req.Amount)), req.From, time.Hour)
	s.jabber.Money.SendInvoice(
		req.From,
		req.ID,
		i,
	)
}
