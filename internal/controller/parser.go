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

	tg "gopkg.in/tucnak/telebot.v2"
)

const (
	messageTypeTransaction = iota
	messageTypeTransactionGroup
	messageTypeUnknown
)

type messageType int

type ParserController interface {
	Parse(m *tg.Message)
	ParseEdited(m *tg.Message)
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

func (c *parserController) Parse(m *tg.Message) {
	t := c.GetTypeFromMessage(m)

	switch t {
	case messageTypeTransaction:
		c.AddTransaction(m)
	case messageTypeTransactionGroup:
		c.AddTransactionGroup(m)
	}
}
func (c *parserController) ParseEdited(m *tg.Message) {
	t := c.GetTypeFromMessage(m)

	switch t {
	case messageTypeTransaction:
		c.UpdateTransaction(m)
	}
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
func (c *parserController) UpdateTransaction(m *tg.Message) {
	amount, description, walletName, err := c.GetParametersFromMessage(&m.Text)
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
	amount, description, walletName, err := c.GetParametersFromMessage(&m.Text)
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
func (c *parserController) errorHandler(m *tg.Message, err error) {
	log.Println(err)
	_, err = c.bot.Send(m.Chat, defines.MessageError)
	if err != nil {
		log.Println(err)
	}
}
func (c *parserController) botRespond(m *tg.Message, msg string) error {
	if _, err := c.bot.Send(m.Chat, msg, tg.ModeMarkdown); err != nil {
		return err
	}
	return nil
}
