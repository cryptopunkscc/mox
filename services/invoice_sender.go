package services

import (
	"time"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/mox/xmpp"
	"github.com/cryptopunkscc/mox/xmpp/money"
)

type InvoiceSender struct {
	xmpp *xmpp.XMPP
	ln   bitcoin.LightningClient
}

func NewInvoiceSender(j *xmpp.XMPP, ln bitcoin.LightningClient) *InvoiceSender {
	s := &InvoiceSender{
		xmpp: j,
		ln:   ln,
	}
	s.xmpp.Money.InvoiceRequestStream.Subscribe(s.onRequest, nil, nil)
	return s
}

func (s *InvoiceSender) onRequest(req *money.InvoiceRequest) {
	i := s.ln.CreateInvoice(bitcoin.Sat(int64(req.Amount)), req.From, time.Hour)
	s.xmpp.Money.SendInvoice(
		req.From,
		req.ID,
		i,
	)
}
