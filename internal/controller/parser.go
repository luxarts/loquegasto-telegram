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

	tg "gopkg.in/telebot.v3"
)

const (
	messageTypeAddTransaction = iota
	messageTypeUnknown
)

type messageType int

type ParserController interface {
	Parse(ctx tg.Context) error
	ParseEdited(ctx tg.Context) error
	GetTypeFromMessage(msg string) messageType
	AddTransaction(ctx tg.Context) error
	UpdateTransaction(ctx tg.Context) error
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
	t := c.GetTypeFromMessage(ctx.Message().Text)

	switch t {
	case messageTypeAddTransaction:
		c.AddTransaction(ctx)
	}
	return nil
}
func (c *parserController) ParseEdited(ctx tg.Context) error {
	t := c.GetTypeFromMessage(ctx.Message().Text)

	switch t {
	case messageTypeAddTransaction:
		return c.UpdateTransaction(ctx)
	}
	return nil
}
func (c *parserController) GetTypeFromMessage(msg string) messageType {
	// Add transaction check
	r := defines.RegexTransaction.FindStringIndex(msg)
	if r != nil {
		return messageTypeAddTransaction
	}

	return messageTypeUnknown
}
func (c *parserController) UpdateTransaction(ctx tg.Context) error {
	amount, description, walletName, err := c.getParametersFromMessage(&ctx.Message().Text)
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
	amount, description, walletName, err := c.getParametersFromMessage(&ctx.Message().Text)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	if walletName == "" {
		walletName = "efectivo"
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
func (c *parserController) errorHandler(ctx tg.Context, err error) {
	log.Println(err)
	err = ctx.Send(defines.MessageError)
	if err != nil {
		log.Println(err)
	}
}
func (c *parserController) botRespond(ctx tg.Context, msg string) error {
	if err := ctx.Send(msg, tg.ModeMarkdown); err != nil {
		return err
	}
	return nil
}
