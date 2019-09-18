package chatbot

import (
	"fmt"

	"github.com/cryptopunkscc/go-rx"
	"github.com/cryptopunkscc/go-xmppc/components/roster"
)

func (bot *ChatBot) cmdRoster(args []string) string {
	sub := "list"
	res := ""
	if len(args) > 0 {
		sub = args[0]
	}

	switch sub {
	case "list":
		res = "Your roster:\n"
		ch := make(chan []*roster.RosterItem, 0)

		bot.jabber.Roster.FetchRoster(rx.Pipe(func(list []*roster.RosterItem) {
			ch <- list
		}))

		list := <-ch

		for _, item := range list {
			res = res + fmt.Sprintf("%s <%s> %s\n", item.Name, item.JID, item.Subscription)
		}
	case "add":
		if len(args) < 3 {
			return "roster add <jid> <name>"
		}
		bot.jabber.Roster.Add(args[1], args[2])
		return "Added."
	case "remove":
		if len(args) < 2 {
			return "roster remove <jid>"
		}
		bot.jabber.Roster.Remove(args[1])
		return "Removed."
	default:
		return fmt.Sprintf("Unknown command: %s", sub)
	}

	return fmt.Sprintf(res)
}
