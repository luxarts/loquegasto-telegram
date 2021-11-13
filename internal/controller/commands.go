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
	bot       *tg.Bot
	txnSrv    service.TransactionsService
	userSrv   service.UsersService
	walletSrv service.WalletsService
}

func NewCommandsController(bot *tg.Bot, txnSrv service.TransactionsService, usersSrv service.UsersService, walletSrv service.WalletsService) CommandsController {
	return &commandsController{
		bot:       bot,
		txnSrv:    txnSrv,
		userSrv:   usersSrv,
		walletSrv: walletSrv,
	}
}

func (c *commandsController) Start(m *tg.Message) {
	// Create user
	if err := c.userSrv.Create(m.Sender.ID, m.Unixtime, m.Chat.ID); err != nil {
		c.errorHandler(m, err)
		return
	}

	// Create default wallet
	if err := c.walletSrv.Create(m.Sender.ID, "Efectivo", 0.0, m.Unixtime); err != nil {
		c.errorHandler(m, err)
		return
	}

	// Show onboarding message
	if _, err := c.bot.Send(m.Sender, fmt.Sprintf(defines.MessageStart, m.Sender.FirstName), tg.ModeMarkdown); err != nil {
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

	var total int64
	for _, txn := range *transacciones {
		total += int64(txn.Amount * 100)
	}

	_, err = c.bot.Send(m.Sender, fmt.Sprintf(defines.MessageConsumosResponse, "Efectivo", float64(total)/100), tg.ModeMarkdown)
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
