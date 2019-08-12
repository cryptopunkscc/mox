package money

import "encoding/xml"

type xmppInvoice struct {
	XMLName xml.Name `xml:"urn:xmpp:bitcoin invoice"`
	Amount  int      `xml:"amount,omitempty"`
	Method  string   `xml:"method,omitempty"`
	Data    string   `xml:"data,omitempty"`
}
