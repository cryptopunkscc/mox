package money

import (
	"encoding/xml"
	"github.com/cryptopunkscc/go-xmpp"
)

type xmppMoney struct {
	XMLName xml.Name `xml:"urn:xmpp:bitcoin:0 money"`
	Invoice *xmppInvoice
	Request *xmppRequest
}

type xmppInvoice struct {
	XMLName xml.Name `xml:"invoice"`
	Encoded string   `xml:"data,omitempty"`
}

type xmppRequest struct {
	XMLName xml.Name `xml:"request"`
	Amount  int      `xml:"amount,omitempty"`
}

func init() {
	xmpp.AddElement(&xmppMoney{})
}
