package payments

import (
	"fmt"
	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmpp/client"
)

// Check if Component satisfies the Handler interface
var _ xmppc.Handler = &Component{}

// Component provides methods to exchange XMPP packets
type Component struct {
	xmppc.Session
	InvoiceRequestHandler

	invoiceHandlers map[string]InvoiceHandler
}

type InvoiceRequest struct {
	JID       xmpp.JID
	ID        string
	Amount    bitcoin.Amount
	component *Component
}

type Invoice struct {
	JID     xmpp.JID
	ID      string
	Invoice string
}

type InvoiceRequestHandler func(*InvoiceRequest)
type InvoiceHandler func(*Invoice)

// SendInvoice sends an invoice in response to the invoice request
func (req *InvoiceRequest) SendInvoice(invoice string) error {
	return req.component.SendInvoice(&Invoice{
		JID:     req.JID,
		ID:      req.ID,
		Invoice: invoice,
	})
}

// RequestInvoice sends an invoice request to a JID
func (m *Component) RequestInvoice(req *InvoiceRequest, handler InvoiceHandler) error {
	if req.ID == "" {
		return fmt.Errorf("InvoiceRequest missing an ID")
	}
	if !req.JID.Valid() {
		return fmt.Errorf("InvoiceRequest contains invalid JID")
	}
	if req.Amount.Sat() <= 0 {
		return fmt.Errorf("InvoiceRequest contains an invalid amount")
	}
	if handler == nil {
		return fmt.Errorf("no invoice handler provided")
	}

	msg := &xmpp.Message{
		To:   req.JID,
		Type: "normal",
	}
	msg.AddChild(&XMPPBitcoin{
		Request: &XMPPRequest{
			ID:     req.ID,
			Amount: int(req.Amount.Msat()),
		},
	})
	err := m.Write(msg)
	if err != nil {
		return err
	}
	// Register handler under "<barejid> <id>" key to further protect from invoice request hijacking
	id := fmt.Sprintf("%s %s", req.JID.Bare(), req.ID)
	m.addInvoiceHandler(id, handler)
	return nil
}

// SendInvoice sends a LN invoice to a JID
func (m *Component) SendInvoice(invoice *Invoice) error {
	msg := &xmpp.Message{
		To:   invoice.JID,
		Type: "normal",
	}
	msg.AddChild(&XMPPBitcoin{
		Invoice: &XMPPInvoice{
			ID:      invoice.ID,
			Encoded: invoice.Invoice,
		},
	})
	err := m.Write(msg)
	if err != nil {
		return err
	}
	return nil
}

func (m *Component) HandleMessage(msg *xmpp.Message) {
	if !msg.Normal() {
		return
	}
	if _, ok := msg.Child(&XMPPBitcoin{}).(*XMPPBitcoin); ok {
		m.handleBitcoin(msg)
	}
}

func (m *Component) HandleStanza(s xmpp.Stanza) {
	xmppc.HandleStanza(m, s)
}

func (m *Component) Online(s xmppc.Session) {
	m.Session = s
}

func (m *Component) Offline(error) {
	m.Session = nil
}

func (m *Component) handleBitcoinInvoice(i *Invoice) {
	id := fmt.Sprintf("%s %s", i.JID.Bare(), i.ID)
	handler, ok := m.invoiceHandlers[id]
	if !ok {
		// Unrequested invoices are ignored for now
		// TODO: Handle unrequested invoices
		return
	}
	handler(i)
	delete(m.invoiceHandlers, i.ID)
}

func (m *Component) handleBitcoinRequest(r *InvoiceRequest) {
	if m.InvoiceRequestHandler != nil {
		m.InvoiceRequestHandler(r)
	}
}

func (m *Component) handleBitcoin(msg *xmpp.Message) {
	btc := msg.Child(&XMPPBitcoin{}).(*XMPPBitcoin)

	if btc.Request != nil {
		req := &InvoiceRequest{
			JID:       msg.From,
			ID:        btc.Request.ID,
			Amount:    bitcoin.Msat(int64(btc.Request.Amount)),
			component: m,
		}
		m.handleBitcoinRequest(req)
	}
	if btc.Invoice != nil {
		i := &Invoice{
			JID:     msg.From,
			ID:      btc.Invoice.ID,
			Invoice: btc.Invoice.Encoded,
		}
		m.handleBitcoinInvoice(i)
	}
}

func (m *Component) addInvoiceHandler(id string, handler InvoiceHandler) {
	if m.invoiceHandlers == nil {
		m.invoiceHandlers = make(map[string]InvoiceHandler)
	}
	m.invoiceHandlers[id] = handler
}
