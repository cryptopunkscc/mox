package chatbot

import (
	"fmt"
	"regexp"
	"strings"

	bitcoin "github.com/cryptopunkscc/go-bitcoin"

	"github.com/cryptopunkscc/mox/chatbot/acl"
	"github.com/cryptopunkscc/mox/services"

	"github.com/cryptopunkscc/go-xmppc/components/chat"
	"github.com/cryptopunkscc/mox/xmpp"
)

type command func([]string) string

type ChatBot struct {
	config         *Config
	acl            *acl.ACL
	xmpp           *xmpp.XMPP
	ln             bitcoin.LightningClient
	commands       map[string]command
	balanceChecker *services.BalanceChecker
	moneySender    *services.MoneySender
	invoiceDecoder *services.InvoiceDecoder
}

func New(cfg *Config, j *xmpp.XMPP, ln bitcoin.LightningClient, bc *services.BalanceChecker, ms *services.MoneySender, id *services.InvoiceDecoder) *ChatBot {
	bot := &ChatBot{
		config:         cfg,
		acl:            acl.New(),
		xmpp:           j,
		ln:             ln,
		balanceChecker: bc,
		moneySender:    ms,
		invoiceDecoder: id,
		commands:       make(map[string]command),
	}

	if bot.config == nil {
		bot.config = &Config{}
	}

	bot.acl.Set(bot.config.AdminJID, acl.Permissions{Access: true})
	bot.xmpp.Chat.MessageStream.Subscribe(bot.onMessage, nil, nil)
	bot.commands["balance"] = bot.cmdBalance
	bot.commands["send"] = bot.cmdSend
	bot.commands["decode"] = bot.cmdDecode
	bot.commands["help"] = bot.cmdHelp
	bot.commands["status"] = bot.cmdStatus
	bot.commands["issue"] = bot.cmdIssue
	bot.commands["pay"] = bot.cmdPay
	bot.commands["allow"] = bot.cmdAllow
	bot.commands["roster"] = bot.cmdRoster
	bot.commands["subscribe"] = bot.cmdSubscribe

	ln.SetInvoiceHandler(bot.handleInvoice)
	return bot
}

func (bot *ChatBot) handleInvoice(i *bitcoin.Invoice) {
	if i.State == 1 {
		msg := fmt.Sprintf("Received %d SAT!", i.Amount.Sat())
		bot.xmpp.Chat.SendMessage(bot.config.AdminJID, msg)
	}
}

func (bot *ChatBot) onMessage(msg *chat.Message) {
	parts := regexp.MustCompile("'.+'|\".+\"|\\S+").FindAllString(msg.Body, -1)

	if len(parts) == 0 {
		return
	}

	cmd := strings.ToLower(parts[0])
	args := parts[1:len(parts)]

	p := bot.acl.Get(msg.From.Bare().String())

	if !p.Access {
		bot.xmpp.Chat.SendMessage(msg.From.String(), "Permission denied")
		return
	}

	if fn, ok := bot.commands[cmd]; ok {
		res := fn(args)
		if res != "" {
			bot.xmpp.Chat.SendMessage(msg.From.String(), res)
		}
	} else {
		bot.xmpp.Chat.SendMessage(msg.From.String(), "Unknown command")
	}
}
