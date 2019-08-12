package chatbot

import (
	"fmt"

	"github.com/cryptopunkscc/mox/chatbot/acl"
)

func (bot *ChatBot) cmdAllow(args []string) string {
	if len(args) < 1 {
		return "allow <jid> - allow jid to control this bot"
	}
	bot.acl.Set(args[0], acl.Permissions{Access: true})
	return fmt.Sprintf("Access granted!")
}
