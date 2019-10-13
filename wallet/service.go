package wallet

import (
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnd"
	"time"
)

type Service struct {
	backend bitcoin.LightningClient
}

func New(cfg *Config) (*Service, error) {
	var err error
	w := &Service{}

	w.backend, err = lnd.Connect(&lnd.Config{
		Host:         cfg.getHost(),
		Port:         cfg.getPort(),
		MacaroonPath: cfg.getMacaroon(),
		CertPath:     cfg.getCert(),
	})
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (srv *Service) Balance() bitcoin.Amount {
	return bitcoin.Sat(int64(srv.backend.Balance()))
}

func (srv *Service) Decode(invoice string) *bitcoin.Invoice {
	return srv.backend.DecodeInvoice(invoice)
}

func (srv *Service) IssueInvoice(amount bitcoin.Amount, memo string, validity time.Duration) *bitcoin.Invoice {
	return srv.backend.CreateInvoice(amount, memo, validity)
}

func (srv *Service) PayInvoice(invoice string) error {
	return srv.backend.PayInvoice(invoice)
}
