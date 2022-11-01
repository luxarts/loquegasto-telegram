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

	tg "gopkg.in/telebot.v3"
)

type CommandsController interface {
	Start(ctx tg.Context) error
	Help(ctx tg.Context) error
	Ping(ctx tg.Context) error
	GetWallets(ctx tg.Context) error
	CreateWallet(ctx tg.Context) error
	AddTransaction(ctx tg.Context) error
}

type commandsController struct {
	bot       *tg.Bot
	txnSrv    service.TransactionsService
	userSrv   service.UsersService
	walletSrv service.WalletsService
	oAuthSrv  service.OAuthService
}

func NewCommandsController(bot *tg.Bot, txnSrv service.TransactionsService, usersSrv service.UsersService, walletSrv service.WalletsService, oAuthSrv service.OAuthService) CommandsController {
	return &commandsController{
		bot:       bot,
		txnSrv:    txnSrv,
		userSrv:   usersSrv,
		walletSrv: walletSrv,
		oAuthSrv:  oAuthSrv,
	}
}

func (c *commandsController) Start(ctx tg.Context) error {
	if ctx.Chat().Type == tg.ChatPrivate {
		return c.startPrivate(ctx)
	} else if ctx.Chat().Type == tg.ChatGroup {
		return c.startGroup(ctx)
	}
	return nil
}
func (c *commandsController) startPrivate(ctx tg.Context) error {
	loginURL := c.oAuthSrv.GetLoginURL(ctx.Sender().ID)

	// Create login button
	selector := c.bot.NewMarkup()
	urlBtn := selector.URL("Iniciar sesión con Google", loginURL)
	selector.Inline(
		selector.Row(urlBtn),
	)

	// Show onboarding message
	if err := ctx.Send(fmt.Sprintf(defines.MessageStart, ctx.Sender().FirstName), tg.ModeMarkdown, selector); err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	return nil
}
func (c *commandsController) startGroup(ctx tg.Context) error {
	// Show onboarding message
	c.botRespond(ctx, fmt.Sprintf("@%s registrado.", ctx.Sender().Username))
	return nil
}

func (c *commandsController) Help(ctx tg.Context) error {
	_, err := c.bot.Send(ctx.Recipient(), defines.MessageHelp, tg.ModeMarkdown)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	return nil
}
func (c *commandsController) Ping(ctx tg.Context) error {
	_, err := c.bot.Send(ctx.Recipient(), "pong")
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	return nil
}
func (c *commandsController) GetWallets(ctx tg.Context) error {
	wallets, err := c.walletSrv.GetAll(ctx.Sender().ID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Build response
	response := fmt.Sprintf("*Billeteras:*")
	for _, w := range *wallets {
		response = fmt.Sprintf("%s\n%s: $%.2f", response, w.Name, w.Balance)
	}

	c.botRespond(ctx, response)
	return nil
}
func (c *commandsController) CreateWallet(ctx tg.Context) error {
	timestamp := time.Unix(ctx.Message().Unixtime, 0)
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: ctx.Sender().ID,
	})

	name, balance, err := c.getWalletNameAndBalance(ctx.Message().Payload)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	wallet, err := c.walletSrv.Create(ctx.Sender().ID, name, balance, &timestamp, token)
	if err, isError := err.(*jsend.Body); isError && err != nil {
		c.errorHandlerResponse(ctx, err)
		return err
	}

	response := fmt.Sprintf(defines.MessageCreateWallet, wallet.Name)
	c.botRespond(ctx, response)
	return nil
}
func (c *commandsController) AddTransaction(ctx tg.Context) error {
	payload := domain.CommandTransactionPayload{}
	if err := payload.Parse(ctx.Message().Payload); err != nil {
		c.errorHandlerResponse(ctx, err)
		return err
	}

	response := fmt.Sprintf(defines.MessageAddPaymentResponse, payload.Description, payload.Amount)
	c.botRespond(ctx, response)
	return nil
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
func (c *commandsController) errorHandler(ctx tg.Context, err error) {
	log.Println(err)
	_, err = c.bot.Send(ctx.Recipient(), defines.MessageError, tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *commandsController) errorHandlerResponse(ctx tg.Context, err error) {
	log.Println(err)
	_, err = c.bot.Send(ctx.Recipient(), fmt.Sprintf(defines.MessageErrorResponse, err.Error()), tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *commandsController) botRespond(ctx tg.Context, msg string) {
	if _, err := c.bot.Send(ctx.Recipient(), msg, tg.ModeMarkdown); err != nil {
		c.errorHandler(ctx, err)
	}
}
