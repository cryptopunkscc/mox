package payments

import (
	"encoding/xml"
	"github.com/cryptopunkscc/go-xmpp"
)

// XMPPBitcoin represents an XMPP packet containing all bitcoin-related requests
type XMPPBitcoin struct {
	XMLName xml.Name `xml:"urn:xmpp:bitcoin:0 bitcoin"`
	Invoice *XMPPInvoice
	Request *XMPPRequest
}

// XMPPInvoice represents a lightning network invoice
type XMPPInvoice struct {
	XMLName xml.Name `xml:"invoice"`

	// ID is the ID of the invoice request to which this invoice is a response
	ID string `xml:"id"`

	// Encoded is a BOLT#11 encoded lightning invoice
	Encoded string `xml:"data,omitempty"`
}

// XMPPRequest represents a request for an LN invoice
type XMPPRequest struct {
	XMLName xml.Name `xml:"request"`

	// ID is a unique ID of the request that should be included in invoice response
	ID string `xml:"id"`

	// Amount in msat
	Amount int `xml:"amount,omitempty"`
}

func initXMPP() {
	xmpp.AddElement(&XMPPBitcoin{})
}
