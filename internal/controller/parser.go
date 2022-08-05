package controller

import (
	"errors"
	"fmt"
	"log"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/service"
	"strconv"
	"strings"

	tg "gopkg.in/tucnak/telebot.v2"
)

const (
	messageTypeAddTransaction = iota
	messageTypeUnknown
)

type messageType int

type ParserController interface {
	Parse(m *tg.Message)
	ParseEdited(m *tg.Message)
	GetTypeFromMessage(msg string) messageType
	AddTransaction(m *tg.Message)
	UpdateTransaction(m *tg.Message)
}
type parserController struct {
	bot       *tg.Bot
	txnSrv    service.TransactionsService
	walletSrv service.WalletsService
}

func NewParserController(bot *tg.Bot, txnSrv service.TransactionsService, walletSrv service.WalletsService) ParserController {
	return &parserController{
		bot:       bot,
		txnSrv:    txnSrv,
		walletSrv: walletSrv,
	}
}

func (c *parserController) Parse(m *tg.Message) {
	t := c.GetTypeFromMessage(m.Text)

	switch t {
	case messageTypeAddTransaction:
		c.AddTransaction(m)
	}
}
func (c *parserController) ParseEdited(m *tg.Message) {
	t := c.GetTypeFromMessage(m.Text)

	switch t {
	case messageTypeAddTransaction:
		c.UpdateTransaction(m)
	}
}
func (c *parserController) GetTypeFromMessage(msg string) messageType {
	// Add transaction check
	r := defines.RegexTransaction.FindStringIndex(msg)
	if r != nil {
		return messageTypeAddTransaction
	}

	return messageTypeUnknown
}
func (c *parserController) UpdateTransaction(m *tg.Message) {
	amount, description, walletName, err := c.getParametersFromMessage(&m.Text)
	if err != nil {
		c.errorHandler(m, err)
		return
	}

	if walletName == "" {
		walletName = "efectivo"
	}

	wallet, err := c.walletSrv.GetByName(walletName, m.Sender.ID)
	if err == repository.ErrNotFound {
		c.botRespond(m, defines.MessageErrorWalletNotFound)
		return
	}
	if err != nil {
		c.errorHandler(m, err)
		return
	}

	err = c.txnSrv.UpdateTransaction(m.Sender.ID, m.ID, amount, description, wallet.ID)
	if err != nil {
		c.errorHandler(m, err)
		return
	}

	msg := fmt.Sprintf(defines.MesssageUpdatePaymentResponse)

	// Respond to the user
	_, err = c.bot.Send(m.Sender,
		msg,
		&tg.SendOptions{
			ReplyTo: m,
		},
		tg.ModeMarkdown,
	)
	if err != nil {
		c.errorHandler(m, err)
	}
}
func (c *parserController) AddTransaction(m *tg.Message) {
	amount, description, walletName, err := c.getParametersFromMessage(&m.Text)
	if err != nil {
		c.errorHandler(m, err)
		return
	}
	if walletName == "" {
		walletName = "efectivo"
	}
	wallet, err := c.walletSrv.GetByName(walletName, m.Sender.ID)
	if err == repository.ErrNotFound {
		c.botRespond(m, defines.MessageErrorWalletNotFound)
		return
	}
	if err != nil {
		c.errorHandler(m, err)
		return
	}
	err = c.txnSrv.AddTransaction(m.Sender.ID, m.ID, amount, description, wallet.ID, m.Unixtime)
	if err != nil {
		c.errorHandler(m, err)
		return
	}
	var msg string
	if amount > 0 {
		msg = fmt.Sprintf(defines.MessageAddPaymentResponseWithWallet, description, amount, walletName)
	} else {
		msg = fmt.Sprintf(defines.MessageAddMoneyResponse, description, amount, walletName)
	}
	// Respond to the user
	_, err = c.bot.Send(m.Sender,
		msg,
		&tg.SendOptions{
			ReplyTo: m,
		},
		tg.ModeMarkdown,
	)
	if err != nil {
		c.errorHandler(m, err)
	}
}

func (c *parserController) getParametersFromMessage(msg *string) (amount float64, description, walletName string, err error) {
	// Search for amount and description
	result := defines.RegexTransaction.FindAllStringSubmatch(*msg, -1)

	// Validate results
	if len(result) != 1 || len(result[0]) != 4 {
		err = errors.New("invalid syntax")
		return
	}

	// Amount capture group 1
	amountStr := result[0][1]

	// Parse decimal as dot for internal usage and colon for response
	amountStr = strings.Replace(amountStr, ",", ".", 1)
	amount, err = strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return
	}

	// Description capture group 2
	description = result[0][2]

	// Wallet will be empty if capture group 3 isn't set
	walletName = result[0][3]

	return
}
func (c *parserController) errorHandler(m *tg.Message, err error) {
	log.Println(err)
	_, err = c.bot.Send(m.Sender, defines.MessageError)
	if err != nil {
		log.Println(err)
	}
}
func (c *parserController) botRespond(m *tg.Message, msg string) error {
	if _, err := c.bot.Send(m.Sender, msg, tg.ModeMarkdown); err != nil {
		return err
	}
	return nil
}
