package controller

import (
	"errors"
	"fmt"
	tg "gopkg.in/telebot.v3"
	"log"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/service"
	"strconv"
	"strings"
)

const (
	messageTypeTransaction = iota
	messageTypeUnknown
)

type messageType int

type EventsController interface {
	Parse(ctx tg.Context) error
	Process(ctx tg.Context) error
	GetTypeFromMessage(m *tg.Message) messageType
	//ParseEdited(ctx tg.Context) error
}
type eventsController struct {
	bot          *tg.Bot
	txnSvc       service.TransactionsService
	txnStatusSvc service.TransactionStatusService
	walletsSvc   service.WalletsService
	catSvc       service.CategoriesService
}

func NewEventsController(bot *tg.Bot, txnSvc service.TransactionsService, txnStatusSvc service.TransactionStatusService, walletsSvc service.WalletsService, catSvc service.CategoriesService) EventsController {
	return &eventsController{
		bot:          bot,
		txnSvc:       txnSvc,
		txnStatusSvc: txnStatusSvc,
		walletsSvc:   walletsSvc,
		catSvc:       catSvc,
	}
}

func (c *eventsController) Parse(ctx tg.Context) error {
	t := c.GetTypeFromMessage(ctx.Message())
	var err error

	switch t {
	case messageTypeTransaction:
		err = c.beginTransaction(ctx)
	}
	return err
}
func (c *eventsController) Process(ctx tg.Context) error {
	userID := ctx.Sender().ID
	txnStatus, err := c.txnStatusSvc.GetByUserID(userID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	switch txnStatus.Status {
	case defines.StatusWalletSelection:
		err = c.walletSelection(ctx, txnStatus)
	case defines.StatusCategorySelection:
		err = c.categorySelection(ctx, txnStatus)
	}

	return err
}

/*func (c *eventsController) ParseEdited(ctx tg.Context) error {
	t := c.getTypeFromMessage(ctx.Message())
	var err error

	switch t {
	case messageTypeTransaction:
		//err = c.UpdateTransaction(ctx)
	}
	return err
}*/
/*
	func (c *eventsController) UpdateTransaction(ctx tg.Context) error {
		amount, description, err := c.GetParametersFromMessage(ctx.Message().Text)
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
*/

func (c *eventsController) beginTransaction(ctx tg.Context) error {
	userID := ctx.Sender().ID

	amount, description, err := c.getParametersFromMessage(ctx.Message().Text)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Create and set status to next step: wallet selection
	err = c.txnStatusSvc.Create(userID, amount, description, ctx.Message().Time(), ctx.Message().ID, defines.StatusWalletSelection)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	wallets, err := c.walletsSvc.GetAll(userID)

	kb := buildWalletsKeyboard(wallets)

	// Respond to the user
	err = ctx.Send(
		"¿Con qué billetera?",
		&tg.SendOptions{
			ReplyTo: ctx.Message(),
		},
		tg.ModeMarkdown,
		kb,
	)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	return nil
}
func (c *eventsController) walletSelection(ctx tg.Context, txnStatus *domain.TransactionStatusDTO) error {
	categories, err := c.catSvc.GetAll(txnStatus.Data.UserID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Update and change status to next step: category selection
	walletID, err := strconv.Atoi(strings.Replace(ctx.Callback().Data, "\f", "", 1))
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	txnStatus.Data.WalletID = walletID
	txnStatus.Status = defines.StatusCategorySelection
	err = c.txnStatusSvc.UpdateByUserID(txnStatus)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	kb := buildCategoriesKeyboard(categories)

	// Respond to the user
	err = ctx.EditOrSend(
		"¿De cuál categoría?",
		&tg.SendOptions{
			ReplyTo: ctx.Message(),
		},
		tg.ModeMarkdown,
		kb,
	)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	return nil
}
func (c *eventsController) categorySelection(ctx tg.Context, txnStatus *domain.TransactionStatusDTO) error {
	catID, err := strconv.Atoi(strings.Replace(ctx.Callback().Data, "\f", "", 1))
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	txnStatus.Data.CategoryID = catID

	cat, err := c.catSvc.GetByID(txnStatus.Data.CategoryID, txnStatus.Data.UserID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	wal, err := c.walletsSvc.GetByID(txnStatus.Data.WalletID, txnStatus.Data.UserID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Create transaction
	err = c.txnSvc.AddTransaction(
		txnStatus.Data.UserID,
		txnStatus.Data.MsgID,
		txnStatus.Data.Amount,
		txnStatus.Data.Description,
		txnStatus.Data.WalletID,
		txnStatus.Data.CategoryID,
		txnStatus.Data.CreatedAt,
	)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Delete status
	err = c.txnStatusSvc.DeleteByUserID(txnStatus.Data.UserID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	var msg string

	if txnStatus.Data.Amount > 0 {
		msg = fmt.Sprintf(defines.MessageAddPaymentResponse,
			txnStatus.Data.Description,
			cat.Name,
			txnStatus.Data.Amount,
			wal.Name,
		)
	} else {
		msg = fmt.Sprintf(defines.MessageAddMoneyResponse,
			txnStatus.Data.Description,
			cat.Name,
			txnStatus.Data.Amount,
			wal.Name,
		)
	}

	err = ctx.Edit(
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

// Utils
func (c *eventsController) getParametersFromMessage(msg string) (amount float64, description string, err error) {
	// Search for amount and description
	result := defines.RegexTransaction.FindAllStringSubmatch(msg, -1)

	// Validate results
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

	return
}
func (c *eventsController) GetTypeFromMessage(m *tg.Message) messageType {
	// Add payment check
	if !m.FromGroup() {
		r := defines.RegexTransaction.FindStringIndex(m.Text)
		if r != nil {
			return messageTypeTransaction
		}
	}

	return messageTypeUnknown
}

func (c *eventsController) errorHandler(ctx tg.Context, err error) {
	log.Println(err)
	_, err = c.bot.Send(ctx.Recipient(), defines.MessageError, tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *eventsController) errorHandlerResponse(ctx tg.Context, err error) {
	log.Println(err)
	_, err = c.bot.Send(ctx.Recipient(), fmt.Sprintf(defines.MessageErrorResponse, err.Error()), tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *eventsController) botRespond(ctx tg.Context, msg string) {
	if err := ctx.Send(msg, tg.ModeMarkdown); err != nil {
		c.errorHandler(ctx, err)
	}
}

func buildWalletsKeyboard(wallets *[]domain.WalletDTO) *tg.ReplyMarkup {
	kb := &tg.ReplyMarkup{}

	var btns []tg.Btn

	for _, w := range *wallets {
		btns = append(btns, kb.Data(w.Name, strconv.Itoa(w.ID)))
	}
	rows := kb.Split(2, btns)
	kb.Inline(rows...)

	return kb
}
func buildCategoriesKeyboard(c *[]domain.CategoryDTO) *tg.ReplyMarkup {
	kb := &tg.ReplyMarkup{}

	var btns []tg.Btn

	for _, c := range *c {
		btns = append(btns, kb.Data(c.Emoji+" "+c.Name, strconv.Itoa(c.ID)))
	}
	rows := kb.Split(2, btns)
	kb.Inline(rows...)

	return kb
}
