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
	"time"

	tg "gopkg.in/telebot.v3"
)

const (
	messageTypeTransaction = iota
	messageTypeTransactionGroup
	messageTypeUnknown
)

type messageType int

type ParserController interface {
	Parse(ctx tg.Context) error
	ParseEdited(ctx tg.Context) error
	GetTypeFromMessage(m *tg.Message) messageType
	AddTransaction(m *tg.Message)
	AddTransactionGroup(m *tg.Message)
	UpdateTransaction(m *tg.Message)
	GetParametersFromMessage(msg *string) (amount float64, description, walletName string, err error)
}
type parserController struct {
	bot       *tg.Bot
	txnSrv    service.TransactionsService
	walletSrv service.WalletsService
	sheetsSrv service.SheetsService
}

func NewParserController(bot *tg.Bot, txnSrv service.TransactionsService, walletSrv service.WalletsService, sheetsSrv service.SheetsService) ParserController {
	return &parserController{
		bot:       bot,
		txnSrv:    txnSrv,
		walletSrv: walletSrv,
		sheetsSrv: sheetsSrv,
	}
}

func (c *parserController) Parse(ctx tg.Context) error {
	t := c.GetTypeFromMessage(ctx.Message().Text)

	switch t {
	case messageTypeTransaction:
		c.AddTransaction(ctx.Message())
	case messageTypeTransactionGroup:
		c.AddTransactionGroup(ctx.Message())
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
	case messageTypeTransaction:
		c.UpdateTransaction(ctx.Message())
	}
	return nil
}

func (c *parserController) GetTypeFromMessage(m *tg.Message) messageType {
	// Add payment check
	if !m.FromGroup() {
		r := defines.RegexTransaction.FindStringIndex(m.Text)
		if r != nil {
			return messageTypeTransaction
		}
	} else {
		// Add group transaction check
		r := defines.RegexTransactionGroup.FindStringIndex(m.Text)
		if r != nil {
			return messageTypeTransactionGroup
		}
	}

	return messageTypeUnknown
}
func (c *parserController) UpdateTransaction(ctx tg.Context) error {
	amount, description, walletName, err := c.GetParametersFromMessage(&ctx.Message().Text)
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
	amount, description, walletName, err := c.GetParametersFromMessage(&ctx.Message().Text)
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
func (c *parserController) AddTransactionGroup(m *tg.Message) {
	amount, description, payerName, err := c.GetAddTransactionData(m)
	if err != nil {
		c.errorHandler(m, err)
	}

	resp, err := c.sheetsSrv.AddRow(time.Unix(m.Unixtime, 0), description, amount, payerName)
	if err != nil {
		c.errorHandler(m, err)
	}
	respRange := strings.Replace(resp.Updates.UpdatedRange, "Gastos!", "", 1)
	msg := "Anotado -> [link](https://docs.google.com/spreadsheets/d/" + c.sheetsSrv.GetSpreadsheetID() + "/edit#gid=0&range=" + respRange + ")"
	_, err = c.bot.Send(m.Chat, msg, tg.ModeMarkdown, tg.NoPreview)

	if err != nil {
		c.errorHandler(m, err)
	}
}
func (c *parserController) GetParametersFromMessage(msg *string) (amount float64, description, walletName string, err error) {
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
