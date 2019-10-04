package roster

import "github.com/cryptopunkscc/go-xmpp"

type View struct {
}

type Contact struct {
	JID  xmpp.JID
	Name string
}

type Presence struct {
	Online bool
	Status string
	Show   string
}

func (view *View) AddContact(jid xmpp.JID, name string) {

}
