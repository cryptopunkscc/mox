package services

import (
	"github.com/cryptopunkscc/go-bitcoin"
)

type BalanceChecker struct {
	ln bitcoin.LightningClient
}

func NewBalanceChecker(ln bitcoin.LightningClient) *BalanceChecker {
	return &BalanceChecker{
		ln: ln,
	}
}

func (service *BalanceChecker) Check() int {
	return service.ln.Balance()
}
