package controller

import (
	"fmt"
	"log"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/service"

	tg "gopkg.in/tucnak/telebot.v2"
)

type CommandsController interface {
	Start(m *tg.Message)
	Help(m *tg.Message)
	Ping(m *tg.Message)
}

type commandsController struct {
	bot    *tg.Bot
	txnSrv service.TransactionsService
}

func NewCommandsController(bot *tg.Bot, txnSrv service.TransactionsService) CommandsController {
	return &commandsController{
		bot:    bot,
		txnSrv: txnSrv,
	}
}

func (c *commandsController) Start(m *tg.Message) {
	// Response
	_, err := c.bot.Send(m.Sender,
		fmt.Sprintf(defines.MessageStart, m.Sender.FirstName),
	)
	if err != nil {
		c.errorHandler(m, err)
		return
	}
}
func (c *commandsController) Help(m *tg.Message) {
	_, err := c.bot.Send(m.Sender, defines.MessageHelp, tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *commandsController) Ping(m *tg.Message) {
	_, err := c.bot.Send(m.Sender, "pong")
	if err != nil {
		log.Println(err)
	}
}

func (c *commandsController) errorHandler(m *tg.Message, err error) {
	log.Println(err)
	_, err = c.bot.Send(m.Sender, defines.MessageError)
	if err != nil {
		log.Println(err)
	}
}
