package contacts

import (
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmpp/client/components/presence"
	"time"
)

type resource struct {
	JID        xmpp.JID
	Status     string
	Priority   int
	Show       string
	Online     bool
	LastUpdate time.Time
}

type contact struct {
	JID       xmpp.JID
	Name      string
	Remove    bool
	Resources map[string]*resource
}

func (c *contact) Online() bool {
	for _, r := range c.Resources {
		if r.Online {
			return true
		}
	}
	return false
}

func (c *contact) Status() string {
	if best := c.BestResource(); best != nil {
		return best.Status
	}
	return ""
}

func (c *contact) UpdatePresence(update *presence.Update) {
	if c.Resources == nil {
		c.Resources = make(map[string]*resource)
	}
	name := update.JID.Resource()
	c.Resources[name] = &resource{
		JID:        update.JID,
		Status:     update.Status,
		Priority:   update.Priority,
		Show:       update.Show,
		Online:     update.Online,
		LastUpdate: time.Now(),
	}
}

func (r *resource) BetterThan(other *resource) bool {
	if other == nil {
		return true
	}
	if r.Priority == other.Priority {
		return r.LastUpdate.After(other.LastUpdate)
	}
	return r.Priority > other.Priority
}

func (c contact) BestResource() (best *resource) {
	for _, r := range c.Resources {
		if r.BetterThan(best) {
			best = r
		}
	}
	return
}
