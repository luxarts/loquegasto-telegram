package controller

import (
	"fmt"
	"log"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/service"
	"strconv"
	"strings"

	tg "gopkg.in/tucnak/telebot.v2"
)

type CommandsController interface {
	Start(m *tg.Message)
	Help(m *tg.Message)
	Ping(m *tg.Message)
	Total(m *tg.Message)
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
	_, err := c.bot.Send(m.Sender, fmt.Sprintf(defines.MessageStart, m.Sender.FirstName), tg.ModeMarkdown)
	if err != nil {
		c.errorHandler(m, err)
	}
}
func (c *commandsController) Help(m *tg.Message) {
	_, err := c.bot.Send(m.Sender, defines.MessageHelp, tg.ModeMarkdown)
	if err != nil {
		c.errorHandler(m, err)
	}
}
func (c *commandsController) Ping(m *tg.Message) {
	_, err := c.bot.Send(m.Sender, "pong")
	if err != nil {
		c.errorHandler(m, err)
	}
}
func (c *commandsController) Total(m *tg.Message) {
	total, err := c.txnSrv.GetTotal(m.Sender.ID)
	if err != nil {
		c.errorHandler(m, err)
		return
	}

	totalStr := strconv.FormatFloat(total, byte('f'), 2, 64)
	totalStr = strings.Replace(totalStr, ".", ",", 1)

	_, err = c.bot.Send(m.Sender, fmt.Sprintf(defines.MessageTotalResponse, totalStr), tg.ModeMarkdown)
	if err != nil {
		c.errorHandler(m, err)
		return
	}
}

// Utils
func (c *commandsController) errorHandler(m *tg.Message, err error) {
	log.Println(err)
	_, err = c.bot.Send(m.Sender, defines.MessageError, tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
