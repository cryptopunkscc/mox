package chatbot

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (bot *ChatBot) cmdDecode(args []string) string {
	if len(args) < 1 {
		return "decode <paymentRequest>"
	}
	payReq := args[0]
	invoice := bot.invoiceDecoder.Decode(payReq)
	fmt := message.NewPrinter(language.English)
	return fmt.Sprintf("Invoice is for %d SAT, memo:\n%s", invoice.Amount.Sat(), invoice.Description)
}
