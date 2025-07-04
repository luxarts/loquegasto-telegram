package controller

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/service"
	"time"

	tg "gopkg.in/telebot.v3"
)

type CommandsController interface {
	Start(ctx tg.Context) error
	Help(ctx tg.Context) error
	Ping(ctx tg.Context) error
	GetWallets(ctx tg.Context) error
	CreateWallet(ctx tg.Context) error
	CreateCategory(ctx tg.Context) error
	Cancel(ctx tg.Context) error
	Export(ctx tg.Context) error
}

type commandsController struct {
	bot         *tg.Bot
	txnSvc      service.TransactionsService
	userSvc     service.UsersService
	walletSvc   service.WalletsService
	usrStateSvc service.UserStateService
	exporterSvc service.ExporterService
	catSvc      service.CategoriesService
}

func NewCommandsController(bot *tg.Bot, txnSvc service.TransactionsService, usersSvc service.UsersService, walletSvc service.WalletsService, usrStateSvc service.UserStateService, exporterSvc service.ExporterService, catSvc service.CategoriesService) CommandsController {
	return &commandsController{
		bot:         bot,
		txnSvc:      txnSvc,
		userSvc:     usersSvc,
		walletSvc:   walletSvc,
		usrStateSvc: usrStateSvc,
		exporterSvc: exporterSvc,
		catSvc:      catSvc,
	}
}

func (c *commandsController) Start(ctx tg.Context) error {
	var err error

	switch ctx.Chat().Type {
	case tg.ChatPrivate:
		err = c.startPrivate(ctx)
	}
	if err != nil {
		c.errorHandler(ctx, err)
	}
	return err
}
func (c *commandsController) startPrivate(ctx tg.Context) error {
	userID := ctx.Sender().ID

	err := c.userSvc.Create(userID)
	if err != nil {
		return err
	}

	// Show onboarding message
	return ctx.Reply(fmt.Sprintf(defines.MessageStart, ctx.Sender().FirstName), tg.ModeMarkdown)
}
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
	token, err := c.userSvc.GetToken(ctx.Sender().ID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	wallets, err := c.walletSvc.GetAll(token)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Build response
	var totalBalance float64
	response := fmt.Sprintf("*Billeteras:*")
	for _, w := range *wallets {
		totalBalance += w.Balance
		response = fmt.Sprintf("%s\n%s: $%s", response, w.Name, formatFloat(w.Balance))
	}
	response = fmt.Sprintf("%s\n\nTotal: $%s", response, formatFloat(totalBalance))

	err = ctx.Send(response, tg.ModeMarkdown)
	if err != nil {
		c.errorHandler(ctx, err)
	}
	return err
}
func (c *commandsController) CreateWallet(ctx tg.Context) error {
	err := c.usrStateSvc.SetState(ctx.Sender().ID, defines.StateCreateWalletWaitingName)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	err = ctx.Send(
		defines.MessageCreateWalletWaitingName,
		tg.ModeMarkdown,
	)

	if err != nil {
		c.errorHandler(ctx, err)
	}
	return err
}
func (c *commandsController) CreateCategory(ctx tg.Context) error {
	err := c.usrStateSvc.SetState(ctx.Sender().ID, defines.StateCreateCategoryWaitingName)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	err = ctx.Send(
		defines.MessageCreateCategoryWaitingName,
		tg.ModeMarkdown,
	)

	if err != nil {
		c.errorHandler(ctx, err)
	}
	return err
}
func (c *commandsController) Cancel(ctx tg.Context) error {
	err := c.usrStateSvc.DeleteByUserID(ctx.Sender().ID)
	if err != nil {
		c.errorHandler(ctx, err)
	}

	err = ctx.Send(defines.MessageCancel, tg.ModeMarkdown)
	if err != nil {
		c.errorHandler(ctx, err)
	}
	return err
}
func (c *commandsController) Export(ctx tg.Context) error {
	userID := ctx.Sender().ID
	token, err := c.userSvc.GetToken(userID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	payload := ctx.Message().Payload
	log.Printf("payload: %s\n", payload)

	var now *time.Time
	var from *time.Time

	if payload != "*" {
		to := ctx.Message().Time()
		now = &to
		fr := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, to.Location())
		from = &fr
	}

	txns, err := c.txnSvc.GetAll(userID, from, now)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	err = c.exporterSvc.Create(userID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Delete file always at the end
	defer func() {
		err := c.exporterSvc.Delete(userID)
		if err != nil {
			c.errorHandler(ctx, err)
		}
	}()

	for _, txn := range *txns {
		wal, err := c.walletSvc.GetByID(txn.WalletID, token)
		if err != nil {
			c.errorHandler(ctx, err)
			return err
		}
		cat, err := c.catSvc.GetByID(txn.CategoryID, token)
		if err != nil {
			c.errorHandler(ctx, err)
			return err
		}

		fixedCreatedAt := txn.CreatedAt.Add(time.Hour * time.Duration(defines.DefaultUserTimeZone)) // TODO Get from user

		err = c.exporterSvc.AddEntry(
			"",
			txn.Amount,
			txn.Description,
			wal.Name,
			cat.Name,
			&fixedCreatedAt,
			userID,
		)
		if err != nil {
			c.errorHandler(ctx, err)
			return err
		}
	}

	err = ctx.Send(&tg.Document{
		File: tg.File{
			FileLocal: c.exporterSvc.GetFile(userID),
		},
		FileName: fmt.Sprintf("lqg-reporte-%s.csv", ctx.Message().Time().Format("2006-01-02")),
	})
	if err != nil {
		c.errorHandler(ctx, err)
	}

	return err
}

// Utils
func formatFloat(n float64) string {
	return message.NewPrinter(language.Spanish).Sprintf("%.2f", n)
}
func (c *commandsController) errorHandler(ctx tg.Context, err error) {
	log.Printf("Error: %+v\n", err)
	_, err = c.bot.Send(ctx.Recipient(), defines.MessageError, tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *commandsController) errorHandlerResponse(ctx tg.Context, err error) {
	log.Printf("Error: %+v\n", err)
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
