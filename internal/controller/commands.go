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
	Consumos(m *tg.Message)
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
func (c *commandsController) Consumos(m *tg.Message) {
	transacciones, err := c.txnSrv.GetAll(m.Sender.ID)
	if err != nil {
		c.errorHandler(m, err)
		return
	}
	fmt.Println(transacciones)

	/*totalStr := strconv.FormatFloat(total, byte('f'), 2, 64)
	totalStr = strings.Replace(totalStr, ".", ",", 1)
	*/
	_, err = c.bot.Send(m.Sender, fmt.Sprintf(defines.MessageConsumosResponse, "Efectivo", 123.45), tg.ModeMarkdown)
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
