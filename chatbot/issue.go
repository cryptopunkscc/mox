package chatbot

import (
	"strconv"
	"time"

	"github.com/cryptopunkscc/go-bitcoin"
)

func (bot *ChatBot) cmdIssue(args []string) string {
	if len(args) < 1 {
		return "issue <amount> <memo>"
	}
	amount, _ := strconv.ParseInt(args[0], 10, 64)
	var memo string
	if len(args) > 1 {
		memo = args[1]
	}
	i := bot.ln.CreateInvoice(bitcoin.Sat(amount), memo, time.Hour)
	return i.PaymentRequest
}
