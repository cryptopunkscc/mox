package payments

import (
	"fmt"
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	*Component
	LightningClient bitcoin.LightningClient
}

func (srv *Service) Balance() bitcoin.Amount {
	return bitcoin.Sat(int64(srv.LightningClient.Balance()))
}

func (srv *Service) Decode(i string) *bitcoin.Invoice {
	return srv.LightningClient.DecodeInvoice(i)
}

func (srv *Service) IssueInvoice(amount bitcoin.Amount, memo string, validity time.Duration) *bitcoin.Invoice {
	return srv.LightningClient.CreateInvoice(amount, memo, validity)
}

func (srv *Service) PayInvoice(invoice string) error {
	return srv.LightningClient.PayInvoice(invoice)
}

func (srv *Service) SendBitcoin(jid xmpp.JID, amount bitcoin.Amount) error {
	req := &InvoiceRequest{
		JID:    jid,
		Amount: amount,
		ID:     uuid.New().String(),
	}
	res := make(chan error, 0)
	err := srv.Component.RequestInvoice(req, func(i *Invoice) {
		decoded := srv.LightningClient.DecodeInvoice(i.Invoice)
		if decoded.Amount != amount {
			res <- fmt.Errorf("received invoice amount different than requested")
			return
		}
		res <- srv.LightningClient.PayInvoice(i.Invoice)
	})
	if err != nil {
		return err
	}
	return <-res
}
