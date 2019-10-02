package services

import (
	"github.com/cryptopunkscc/go-bitcoin"
	"time"
)

type LightningService struct {
	LNClient bitcoin.LightningClient
}

func (srv *LightningService) Balance() bitcoin.Amount {
	return bitcoin.Sat(int64(srv.LNClient.Balance()))
}

func (srv *LightningService) Decode(i string) *bitcoin.Invoice {
	return srv.LNClient.DecodeInvoice(i)
}

func (srv *LightningService) IssueInvoice(amount bitcoin.Amount, memo string, validity time.Duration) *bitcoin.Invoice {
	return srv.LNClient.CreateInvoice(amount, memo, validity)
}

func (srv *LightningService) PayInvoice(invoice string) error {
	return srv.LNClient.PayInvoice(invoice)
}
