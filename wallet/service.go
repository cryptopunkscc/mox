package wallet

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-bitcoin/lnd"
	"log"
	"time"
)

type Service struct {
	wallet bitcoin.Wallet
}

func New(cfg *Config) (*Service, error) {
	var err error
	srv := &Service{}

	srv.wallet, err = lnd.New(&lnd.Config{
		Host:        cfg.getHost(),
		Port:        cfg.getPort(),
		Macaroon:    cfg.getMacaroon(),
		Certificate: cfg.getCert(),
	})
	if err != nil {
		return nil, err
	}

	network, _ := srv.wallet.Network(context.Background())
	agent, _ := srv.wallet.Agent(context.Background())

	if network != "testnet" {
		return nil, errors.New("only testnet is supported")
	}

	log.Printf("Using bitcoin agent: %s on %s", agent, network)

	return srv, nil
}

func (srv *Service) Balance() bitcoin.Amount {
	amount, _ := srv.wallet.Lightning().Balance(context.Background())
	return amount
}

func (srv *Service) Decode(invoice string) *bitcoin.Invoice {
	decoded := srv.wallet.Lightning().Decode(context.Background(), invoice)
	return decoded
}

func (srv *Service) IssueInvoice(amount bitcoin.Amount, memo string, validFor time.Duration) *bitcoin.Invoice {
	i, _ := srv.wallet.Lightning().Issue(context.Background(), bitcoin.InvoiceRequest{
		Amount:   amount,
		Memo:     memo,
		ValidFor: validFor,
	})
	return i
}

func (srv *Service) PayInvoice(invoice string) error {
	return srv.wallet.Lightning().Pay(context.Background(), invoice)
}
