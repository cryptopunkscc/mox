package chatbot

import "fmt"

func (bot *ChatBot) cmdPay(args []string) string {
	if len(args) < 1 {
		return "pay <invoice>"
	}
	payreq := args[0]
	err := bot.ln.PayInvoice(payreq)
	if err != nil {
		return fmt.Sprintln("Error:", err)
	}
	return "Paid!"
}
