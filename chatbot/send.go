package chatbot

import (
	"fmt"
	"strconv"
)

func (bot *ChatBot) cmdSend(args []string) string {
	if len(args) < 2 {
		return "send <address> <amount>"
	}
	to := args[0]
	amount, _ := strconv.ParseInt(args[1], 10, 64)
	if err := bot.moneySender.Send(to, int(amount)); err != nil {
		return fmt.Sprintln("Error:", err)
	}
	return "Sent!"
}
