package controller

import (
	"errors"
	"fmt"
	tg "gopkg.in/telebot.v3"
	"log"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/service"
	"strconv"
	"strings"
)

const (
	messageTypeAddTransaction = iota
	messageTypeTransactionGroup
	messageTypeUnknown
)

type messageType int

type ParserController interface {
	Parse(ctx tg.Context) error
	ParseEdited(ctx tg.Context) error
	GetTypeFromMessage(ctx *tg.Message) messageType
	AddTransaction(ctx tg.Context) error
	UpdateTransaction(ctx tg.Context) error
	GetParametersFromMessage(msg string) (amount float64, description, walletName string, err error)
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

func (c *parserController) Parse(ctx tg.Context) error {
	t := c.GetTypeFromMessage(ctx.Message())
	var err error

	switch t {
	case messageTypeAddTransaction:
		err = c.AddTransaction(ctx)
	}
	return err
}
func (c *parserController) ParseEdited(ctx tg.Context) error {
	t := c.GetTypeFromMessage(ctx.Message())
	var err error

	switch t {
	case messageTypeAddTransaction:
		err = c.UpdateTransaction(ctx)
	}
	return err
}

func (c *parserController) GetTypeFromMessage(m *tg.Message) messageType {
	// Add payment check
	if !m.FromGroup() {
		r := defines.RegexTransaction.FindStringIndex(m.Payload)
		if r != nil {
			return messageTypeAddTransaction
		}
	} else {
		// Add group transaction check
		r := defines.RegexTransactionGroup.FindStringIndex(m.Payload)
		if r != nil {
			return messageTypeTransactionGroup
		}
	}

	return messageTypeUnknown
}
func (c *parserController) UpdateTransaction(ctx tg.Context) error {
	amount, description, walletName, err := c.GetParametersFromMessage(ctx.Message().Payload)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	if walletName == "" {
		walletName = "efectivo"
	}

	wallet, err := c.walletSrv.GetByName(walletName, ctx.Sender().ID)
	if err == repository.ErrNotFound {
		c.botRespond(ctx, defines.MessageErrorWalletNotFound)
		return err
	}
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	err = c.txnSrv.UpdateTransaction(ctx.Sender().ID, ctx.Message().ID, amount, description, wallet.ID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	msg := fmt.Sprintf(defines.MesssageUpdatePaymentResponse)

	// Respond to the user
	err = ctx.Send(
		msg,
		&tg.SendOptions{
			ReplyTo: ctx.Message(),
		},
		tg.ModeMarkdown,
	)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	return nil
}
func (c *parserController) AddTransaction(ctx tg.Context) error {
	amount, description, walletName, err := c.GetParametersFromMessage(ctx.Message().Payload)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	if walletName == "" {
		walletName = defines.DefaultWalletName
	}
	wallet, err := c.walletSrv.GetByName(walletName, ctx.Message().Sender.ID)
	if err == repository.ErrNotFound {
		c.botRespond(ctx, defines.MessageErrorWalletNotFound)
		return err
	}
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	err = c.txnSrv.AddTransaction(ctx.Sender().ID, ctx.Message().ID, amount, description, wallet.ID, ctx.Message().Unixtime)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	var msg string
	if amount > 0 {
		msg = fmt.Sprintf(defines.MessageAddPaymentResponseWithWallet, description, amount, walletName)
	} else {
		msg = fmt.Sprintf(defines.MessageAddMoneyResponse, description, amount, walletName)
	}
	// Respond to the user
	err = ctx.Send(
		msg,
		&tg.SendOptions{
			ReplyTo: ctx.Message(),
		},
		tg.ModeMarkdown,
	)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	return nil
}
func (c *parserController) GetParametersFromMessage(msg string) (amount float64, description, walletName string, err error) {
	// Search for amount and description
	result := defines.RegexTransaction.FindAllStringSubmatch(msg, -1)

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
func (c *parserController) GetAddTransactionData(m *tg.Message) (amount float64, description string, payerName string, err error) {
	result := defines.RegexTransactionGroup.FindAllStringSubmatch(m.Text, -1)

	if len(result) != 1 || len(result[0]) != 3 {
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

	if len(m.Entities) != 1 {
		err = errors.New("invalid syntax")
		return
	}

	if m.Entities[0].Type == "mention" && m.Entities[0].User == nil {
		// Get username from message
		payerName = m.Text[m.Entities[0].Offset+1:]
	} else {
		// Get name from user
		payerName = m.Entities[0].User.FirstName + " " + m.Entities[0].User.LastName
	}

	return
}

func (c *parserController) errorHandler(ctx tg.Context, err error) {
	log.Println(err)
	_, err = c.bot.Send(ctx.Recipient(), defines.MessageError, tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *parserController) errorHandlerResponse(ctx tg.Context, err error) {
	log.Println(err)
	_, err = c.bot.Send(ctx.Recipient(), fmt.Sprintf(defines.MessageErrorResponse, err.Error()), tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *parserController) botRespond(ctx tg.Context, msg string) {
	if err := ctx.Send(msg, tg.ModeMarkdown); err != nil {
		c.errorHandler(ctx, err)
	}
}
