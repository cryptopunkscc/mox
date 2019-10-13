package payments

import (
	"fmt"
	"log"
	"time"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/mox/wallet"
	"github.com/google/uuid"
)

type Service struct {
	*Component
	Wallet *wallet.Service
}

func (srv *Service) Balance() bitcoin.Amount {
	return srv.Wallet.Balance()
}

func (srv *Service) Decode(i string) *bitcoin.Invoice {
	return srv.Wallet.Decode(i)
}

func (srv *Service) IssueInvoice(amount bitcoin.Amount, memo string, validity time.Duration) *bitcoin.Invoice {
	return srv.Wallet.IssueInvoice(amount, memo, validity)
}

func (srv *Service) PayInvoice(invoice string) error {
	return srv.Wallet.PayInvoice(invoice)
}

func (srv *Service) SendBitcoin(jid xmpp.JID, amount bitcoin.Amount) error {
	if jid.Resource() == "" {
		jid = jid + "/mox"
	}

	req := &InvoiceRequest{
		JID:    jid,
		Amount: amount,
		ID:     uuid.New().String(),
	}
	res := make(chan error, 0)

	log.Printf("SendBitcoin: Requesting invoice for %d SAT from %s...", req.Amount.Sat(), req.JID)
	err := srv.Component.RequestInvoice(req, func(i *Invoice) {
		log.Printf("SendBitcoin: invoice received!")
		decoded := srv.Wallet.Decode(i.Invoice)
		if decoded.Amount != amount {
			res <- fmt.Errorf("received invoice amount different than requested")
			return
		}
		res <- srv.Wallet.PayInvoice(i.Invoice)
	})
	if err != nil {
		return err
	}
	return <-res
}
