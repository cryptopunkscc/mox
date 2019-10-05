package roster

import (
	"errors"
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmppc/components/presence"
	"github.com/cryptopunkscc/go-xmppc/components/roster"
	"sync"
)

type Service struct {
	Roster   *roster.Roster
	Presence *presence.Presence

	contacts map[xmpp.JID]*contact
	mu       sync.Mutex
}

type Contact struct {
	JID    xmpp.JID
	Name   string
	Online bool
	Status string
}

func (srv *Service) AvailableContacts() []Contact {
	res := make([]Contact, 0)
	for _, c := range srv.contacts {
		res = append(res, Contact{
			JID:    c.JID,
			Name:   c.Name,
			Online: c.Online(),
			Status: c.BestResource().Status,
		})
	}
	return res
}

func (srv *Service) AddContact(jid xmpp.JID, name string) error {
	bare := jid.Bare()
	srv.Roster.Add(bare, name)
	srv.Presence.Subscribe(bare)
	return srv.addContact(bare, name)
}

func (srv *Service) RemoveContact(jid xmpp.JID) error {
	srv.Roster.Remove(jid.Bare())
	return srv.removeContact(jid.Bare())
}

func (srv *Service) addContact(jid xmpp.JID, name string) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.contacts == nil {
		srv.contacts = make(map[xmpp.JID]*contact)
	}
	if _, found := srv.contacts[jid]; found {
		return errors.New("contact already added")
	}
	if name == "" {
		name = jid.Local()
	}
	srv.contacts[jid] = &contact{
		JID:  jid,
		Name: name,
	}
	return nil
}

func (srv *Service) removeContact(jid xmpp.JID) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.contacts == nil {
		return errors.New("contact not found")
	}
	if _, found := srv.contacts[jid]; !found {
		return errors.New("contact not found")
	}
	delete(srv.contacts, jid)
	return nil
}

func (srv *Service) Fetch() {
	srv.Roster.FetchRoster(srv.patchRoster)
}

func (srv *Service) UpdatePresence(update *presence.Update) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	bare := update.JID.Bare()
	if srv.contacts == nil {
		srv.contacts = make(map[xmpp.JID]*contact, 0)
	}
	if _, ok := srv.contacts[bare]; !ok {
		srv.contacts[bare] = &contact{
			JID:  bare,
			Name: update.JID.Local(),
		}
	}
	srv.contacts[bare].UpdatePresence(update)
}

func (srv *Service) patchRoster(list []*roster.RosterItem) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.contacts == nil {
		srv.contacts = make(map[xmpp.JID]*contact)
	}
	for _, c := range srv.contacts {
		c.Remove = true
	}
	for _, item := range list {
		if _, found := srv.contacts[item.JID]; !found {
			srv.contacts[item.JID] = &contact{}
		}
		c := srv.contacts[item.JID]
		c.JID = item.JID
		c.Name = item.Name
		c.Remove = false
		srv.contacts[item.JID] = c
	}
	for k, c := range srv.contacts {
		if c.Remove {
			delete(srv.contacts, k)
		}
	}
}
