package controller

import (
	"errors"
	"fmt"
	"github.com/luxarts/jsend-go"
	"log"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/service"
	"loquegasto-telegram/internal/utils/jwt"
	"strconv"
	"strings"
	"time"

	tg "gopkg.in/tucnak/telebot.v2"
)

type CommandsController interface {
	Start(m *tg.Message)
	Help(m *tg.Message)
	Ping(m *tg.Message)
	GetWallets(m *tg.Message)
	CreateWallet(m *tg.Message)
	AddTransaction(m *tg.Message)
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
	if m.Chat.Type == tg.ChatPrivate {
		c.startPrivate(m)
	} else if m.Chat.Type == tg.ChatGroup {
		c.startGroup(m)
	}
}
func (c *commandsController) startPrivate(m *tg.Message) {
	timestamp := time.Unix(m.Unixtime, 0)
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: m.Sender.ID,
	})

	// Create user
	if err := c.userSrv.Create(m.Sender.ID, &timestamp, m.Chat.ID, token); err != nil {
		c.errorHandler(m, err)
		return
	}

	// Create default wallet
	if _, err := c.walletSrv.Create(m.Sender.ID, "Efectivo", 0.0, &timestamp, token); err != nil {
		c.errorHandler(m, err)
		return
	}

	// Show onboarding message
	if _, err := c.bot.Send(m.Sender, fmt.Sprintf(defines.MessageStart, m.Sender.FirstName), tg.ModeMarkdown); err != nil {
		c.errorHandler(m, err)
	}
}
func (c *commandsController) startGroup(m *tg.Message) {
	// Show onboarding message
	c.botRespond(m, fmt.Sprintf("@%s registrado.", m.Sender.Username))
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
func (c *commandsController) GetWallets(m *tg.Message) {
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

	c.botRespond(m, response)
}
func (c *commandsController) CreateWallet(m *tg.Message) {
	timestamp := time.Unix(m.Unixtime, 0)
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: m.Sender.ID,
	})

	name, balance, err := c.getWalletNameAndBalance(m.Payload)
	if err != nil {
		c.errorHandler(m, err)
		return
	}

	wallet, err := c.walletSrv.Create(m.Sender.ID, name, balance, &timestamp, token)
	if err, isError := err.(*jsend.Body); isError && err != nil {
		c.errorHandlerResponse(m, err)
		return
	}

	response := fmt.Sprintf(defines.MessageCreateWallet, wallet.Name)
	c.botRespond(m, response)
}
func (c *commandsController) AddTransaction(m *tg.Message) {
	payload := domain.CommandTransactionPayload{}
	if err := payload.Parse(m.Payload); err != nil {
		c.errorHandlerResponse(m, err)
		return
	}

	response := fmt.Sprintf(defines.MessageAddPaymentResponse, payload.Description, payload.Amount)
	c.botRespond(m, response)
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
	_, err = c.bot.Send(m.Chat, defines.MessageError, tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *commandsController) errorHandlerResponse(m *tg.Message, err error) {
	log.Println(err)
	_, err = c.bot.Send(m.Chat, fmt.Sprintf(defines.MessageErrorResponse, err.Error()), tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *commandsController) botRespond(m *tg.Message, msg string) {
	if _, err := c.bot.Send(m.Chat, msg, tg.ModeMarkdown); err != nil {
		c.errorHandler(m, err)
	}
}
