package chatbot

func (bot *ChatBot) cmdHelp([]string) string {
	return `balance   - check balance
send <address> <amount>   - send money
decode <invoice>   - show details about an invoice`
}
