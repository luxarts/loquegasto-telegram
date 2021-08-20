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

const (
	messageTypePayment = iota
	messageTypeUnknown
)

type messageType int

type ParserController interface {
	Parse(m *tg.Message)
	ParseEdited(m *tg.Message)
	GetTypeFromMessage(msg string) messageType
	AddPayment(m *tg.Message)
	UpdatePayment(m *tg.Message)
}
type parserController struct {
	bot    *tg.Bot
	txnSrv service.TransactionsService
}

func NewParserController(bot *tg.Bot, txnSrv service.TransactionsService) ParserController {
	return &parserController{
		bot:    bot,
		txnSrv: txnSrv,
	}
}

func (c *parserController) Parse(m *tg.Message) {
	t := c.GetTypeFromMessage(m.Text)

	switch t {
	case messageTypePayment:
		c.AddPayment(m)
	}
}
func (c *parserController) ParseEdited(m *tg.Message) {
	t := c.GetTypeFromMessage(m.Text)

	switch t {
	case messageTypePayment:
		c.UpdatePayment(m)
	}
}
func (c *parserController) GetTypeFromMessage(msg string) messageType {
	// Add payment check
	r := defines.RegexPayment.FindStringIndex(msg)
	if r != nil {
		return messageTypePayment
	}

	return messageTypeUnknown
}
func (c *parserController) UpdatePayment(m *tg.Message) {
	amount, description, source, err := c.getParametersFromMessage(&m.Text)
	if err != nil {
		c.errorHandler(m, err)
		return
	}

	err = c.txnSrv.UpdatePayment(m.Sender.ID, m.ID, amount, description, source, m.Unixtime)
	if err != nil {
		c.errorHandler(m, err)
		return
	}

	msg := fmt.Sprintf(defines.MesssagePaymentUpdatedResponse)

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
func (c *parserController) AddPayment(m *tg.Message) {
	amount, description, source, err := c.getParametersFromMessage(&m.Text)
	if err != nil {
		c.errorHandler(m, err)
		return
	}

	err = c.txnSrv.AddPayment(m.Sender.ID, m.ID, amount, description, source, m.Unixtime)
	if err != nil {
		c.errorHandler(m, err)
		return
	}

	msg := fmt.Sprintf(defines.MessagePaymentResponse, description, amount)
	if source != "" {
		msg = fmt.Sprintf(defines.MessagePaymentResponseWithSource, description, amount, source)
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

func (c *parserController) getParametersFromMessage(msg *string) (amount float64, description, source string, err error) {
	// Search for amount and description
	result := defines.RegexPayment.FindAllStringSubmatch(*msg, -1)

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

	// Source will be empty if capture group 3 isn't set
	source = result[0][3]

	return
}
func (c *parserController) errorHandler(m *tg.Message, err error) {
	log.Println(err)
	_, err = c.bot.Send(m.Sender, defines.MessageError)
	if err != nil {
		log.Println(err)
	}
}
