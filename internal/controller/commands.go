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
	Start(c tg.Context) error
	Help(c tg.Context) error
	Ping(c tg.Context) error
	Wallets(c tg.Context) error
	CreateWallet(c tg.Context) error
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
	timestamp := time.Unix(ctx.Message().Unixtime, 0)
	token := jwt.GenerateToken(nil, &jwt.Payload{
		Subject: ctx.Sender().ID,
	})

	// Create user
	if err := c.userSrv.Create(ctx.Sender().ID, &timestamp, ctx.Message().Chat.ID, token); err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Create default wallet
	if _, err := c.walletSrv.Create(ctx.Sender().ID, "Efectivo", 0.0, &timestamp, token); err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Show onboarding message
	if err := c.botRespond(ctx, fmt.Sprintf(defines.MessageStart, ctx.Sender().FirstName)); err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	return nil
}
func (c *commandsController) Help(ctx tg.Context) error {
	err := ctx.Send(defines.MessageHelp, tg.ModeMarkdown)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	return nil
}
func (c *commandsController) Ping(ctx tg.Context) error {
	err := ctx.Send("pong")
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	return nil
}
func (c *commandsController) Wallets(ctx tg.Context) error {
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

	if err := c.botRespond(ctx, response); err != nil {
		c.errorHandler(ctx, err)
		return err
	}
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
	if err := c.botRespond(ctx, response); err != nil {
		c.errorHandler(ctx, err)
		return err
	}
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
	err = ctx.Send(defines.MessageError, tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *commandsController) errorHandlerResponse(ctx tg.Context, err error) {
	log.Println(err)
	err = ctx.Send(fmt.Sprintf(defines.MessageErrorResponse, err.Error()), tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *commandsController) botRespond(ctx tg.Context, msg string) error {
	if err := ctx.Send(msg, tg.ModeMarkdown); err != nil {
		return err
	}
	return nil
}
