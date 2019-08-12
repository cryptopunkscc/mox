package chatbot

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (bot *ChatBot) cmdBalance(args []string) string {
	balance := bot.balanceChecker.Check()
	fmt := message.NewPrinter(language.English)
	return fmt.Sprintf("Your balance is %d SAT", balance)
}
