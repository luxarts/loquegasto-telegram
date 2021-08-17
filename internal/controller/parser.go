package controller

import (
	"errors"
	"fmt"
	"log"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/service"
	"strconv"

	tg "gopkg.in/tucnak/telebot.v2"
)

const (
	messageTypeAddPayment = iota
	messageTypeUnknown
)

type messageType int

type ParserController interface {
	Parse(m *tg.Message)
	GetTypeFromMessage(msg string) messageType
	AddPayment(m *tg.Message)
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
	case messageTypeAddPayment:
		c.AddPayment(m)
	}
}
func (c *parserController) GetTypeFromMessage(msg string) messageType {
	// Add payment check
	r := defines.RegexAddPayment.FindStringIndex(msg)
	if r != nil {
		return messageTypeAddPayment
	}

	return messageTypeUnknown
}
func (c *parserController) AddPayment(m *tg.Message) {
	// Search for amount and description
	result := defines.RegexAddPayment.FindAllStringSubmatch(m.Text, -1)

	// Validate results
	if len(result) != 1 || len(result[0]) < 3 || len(result[0]) > 4 {
		c.errorHandler(m, errors.New("invalid-syntax"))
		return
	}

	// Amount capture group 1
	amount, err := strconv.ParseInt(result[0][1], 10, 64)
	if err != nil {
		c.errorHandler(m, err)
		return
	}
	// Description capture group 2
	description := result[0][2]

	source := ""
	// Check if source is provided (capture group 3)
	if len(result[0]) == 4 {
		source = result[0][3]
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

func (c *parserController) errorHandler(m *tg.Message, err error) {
	log.Println(err)
	_, err = c.bot.Send(m.Sender, defines.MessageError)
	if err != nil {
		log.Println(err)
	}
}
