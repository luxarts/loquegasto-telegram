package controller

import (
	"errors"
	"fmt"
	"github.com/luxarts/jsend-go"
	"log"
	"loquegasto-telegram/internal/defines"
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

func (c *commandsController) Start(ctx tg.Context) error {
	var err error

	switch ctx.Chat().Type {
	case tg.ChatPrivate:
		err = c.startPrivate(ctx)
	case tg.ChatGroup:
		//err = c.startGroup(ctx)
	}
	if err != nil {
		c.errorHandler(ctx, err)
	}
	return err
}
func (c *commandsController) startPrivate(ctx tg.Context) error {
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: ctx.Sender().ID,
	})
	ts := time.Unix(ctx.Message().Unixtime, 0)

	// Create user
	err := c.userSrv.Create(ctx.Sender().ID, &ts, ctx.Chat().ID, token)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Crear wallet
	_, err = c.walletSrv.Create(ctx.Sender().ID, defines.DefaultWalletName, 0, &ts, token)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Show onboarding message
	return ctx.Send(fmt.Sprintf(defines.MessageStart, ctx.Sender().FirstName), tg.ModeMarkdown)
}

/*func (c *commandsController) startGroup(ctx tg.Context) error {
	// Show onboarding message
	c.botRespond(ctx, fmt.Sprintf("@%s registrado.", ctx.Sender().Username))
	return nil
}*/

func (c *commandsController) Help(ctx tg.Context) error {
	err := ctx.Send(defines.MessageHelp, tg.ModeMarkdown)
	if err != nil {
		c.errorHandler(ctx, err)
	}
	return err
}
func (c *commandsController) Ping(ctx tg.Context) error {
	err := ctx.Send("pong")
	if err != nil {
		c.errorHandler(ctx, err)
	}
	return err
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

	err = ctx.Send(response, tg.ModeMarkdown)
	if err != nil {
		c.errorHandler(ctx, err)
	}
	return err
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
	err = ctx.Send(response, tg.ModeMarkdown)
	if err != nil {
		c.errorHandler(ctx, err)
	}
	return err
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
	if err := ctx.Send(msg, tg.ModeMarkdown); err != nil {
		c.errorHandler(ctx, err)
	}
}
