package chatbot

func (bot *ChatBot) cmdSubscribe(args []string) string {
	if len(args) < 1 {
		return "subscribe <jid>"
	}
	to := args[0]
	bot.xmpp.Presence.Subscribe(to)
	return "Subscription request sent!"
}
