package services

import (
	"github.com/cryptopunkscc/go-bitcoin"
)

type InvoiceDecoder struct {
	ln bitcoin.LightningClient
}

func NewInvoiceDecoder(ln bitcoin.LightningClient) *InvoiceDecoder {
	return &InvoiceDecoder{
		ln: ln,
	}
}

func (service *InvoiceDecoder) Decode(i string) *bitcoin.Invoice {
	return service.ln.DecodeInvoice(i)
}
