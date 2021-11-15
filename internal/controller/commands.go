package controller

import (
	"errors"
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
	Wallets(m *tg.Message)
	CreateWallet(m *tg.Message)
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
func (c *commandsController) Wallets(m *tg.Message) {
	wallets, err := c.walletSrv.GetAll(m.Sender.ID)
	if err != nil {
		c.errorHandler(m, err)
		return
	}

	// Build response
	response := fmt.Sprintf("*Billeteras:*")
	for _, w := range *wallets {
		response = fmt.Sprintf("%s\n%s: $%.2f", response, w.Name, w.Balance)
	}

	if err := c.botRespond(m, response); err != nil {
		c.errorHandler(m, err)
		return
	}
}
func (c *commandsController) CreateWallet(m *tg.Message) {
	name, balance, err := c.getWalletNameAndBalance(m.Payload)
	if err != nil {
		c.errorHandler(m, err)
		return
	}
	err = c.walletSrv.Create(m.Sender.ID, name, balance, m.Unixtime)
	if err != nil {
		c.errorHandler(m, err)
		return
	}

	response := fmt.Sprintf(defines.MessageCreateWallet, name)
	if err := c.botRespond(m, response); err != nil {
		c.errorHandler(m, err)
		return
	}
}

// Utils
func (c *commandsController) getWalletNameAndBalance(text string) (name string, balance float64, err error) {
	result := defines.RegexCreateWallet.FindAllStringSubmatch(text, -1)

	// Validate results
	if len(result) != 1 {
		err = errors.New("invalid syntax")
		return
	}

	// Name capture group 1
	name = result[0][1]

	// Balance capture group 2
	balanceStr := result[0][2]

	// Parse decimal as dot for internal usage and colon for response
	balanceStr = strings.Replace(balanceStr, ",", ".", 1)
	balance, err = strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return
	}

	return
}
func (c *commandsController) errorHandler(m *tg.Message, err error) {
	log.Println(err)
	_, err = c.bot.Send(m.Sender, defines.MessageError, tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *commandsController) botRespond(m *tg.Message, msg string) error {
	if _, err := c.bot.Send(m.Sender, msg, tg.ModeMarkdown); err != nil {
		return err
	}
	return nil
}
