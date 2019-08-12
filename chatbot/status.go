package chatbot

func (bot *ChatBot) cmdStatus(args []string) string {
	if len(args) < 1 {
		return "status <newstatus>"
	}
	bot.jabber.Presence.SetStatus(args[0])
	return "Status changed!"
}
