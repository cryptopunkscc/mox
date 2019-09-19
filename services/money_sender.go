package services

import (
	"log"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-rx"
	"github.com/cryptopunkscc/go-xmpp"
)

type MoneySender struct {
	invoiceRequster *InvoiceRequster
	ln              bitcoin.LightningClient
	invoiceStream   rx.SyncStream
}

func NewMoneySender(r *InvoiceRequster, ln bitcoin.LightningClient) *MoneySender {
	s := &MoneySender{
		invoiceRequster: r,
		ln:              ln,
	}
	s.invoiceStream.Subscribe(s.onPaymentRequest, nil, nil)
	return s
}

func (p *MoneySender) onPaymentRequest(req string) {
	log.Println("Received payment request:", req)
	err := p.ln.PayInvoice(req)
	if err != nil {
		log.Println("Error paying invoice:", err)
	}
}

func (p *MoneySender) Send(to string, amount int) error {
	jid := xmpp.JID(to)

	if jid.Resource() == "" {
		jid = jid.Bare() + "/mox"
	}

	p.invoiceRequster.Request(jid.String(), amount, &p.invoiceStream)

	return nil
}
