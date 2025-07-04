package controller

import (
	"errors"
	"fmt"
	tg "gopkg.in/telebot.v3"
	"log"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/domain"
	"loquegasto-telegram/internal/service"
	"loquegasto-telegram/internal/utils/maptostruct"
	"strconv"
	"strings"
	"time"
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
	bot         *tg.Bot
	txnSvc      service.TransactionsService
	usrStateSvc service.UserStateService
	walletsSvc  service.WalletsService
	catSvc      service.CategoriesService
	usrSvc      service.UsersService
}

func NewEventsController(bot *tg.Bot, txnSvc service.TransactionsService, usrStateSvc service.UserStateService, walletsSvc service.WalletsService, catSvc service.CategoriesService, usrSvc service.UsersService) EventsController {
	return &eventsController{
		bot:         bot,
		txnSvc:      txnSvc,
		usrStateSvc: usrStateSvc,
		walletsSvc:  walletsSvc,
		catSvc:      catSvc,
		usrSvc:      usrSvc,
	}
}

func (c *eventsController) Parse(ctx tg.Context) error {
	usrState, err := c.usrStateSvc.GetByUserID(ctx.Sender().ID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// If user has a state
	if usrState != nil {
		err = c.processState(ctx, usrState)
		if err != nil {
			c.errorHandler(ctx, err)
			return err
		}
		return nil
	}

	t := c.GetTypeFromMessage(ctx.Message())

	switch t {
	case messageTypeTransaction:
		err = c.beginCreateTransaction(ctx)
	default:
		return nil
	}
	return err
}
func (c *eventsController) Process(ctx tg.Context) error {
	userID := ctx.Sender().ID
	txnStatus, err := c.usrStateSvc.GetByUserID(userID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	switch txnStatus.State {
	case defines.StateCreateTransactionSelectingWallet:
		err = c.walletSelection(ctx, txnStatus)
	case defines.StateCreateTransactionSelectingCategory:
		err = c.categorySelection(ctx, txnStatus)
	}

	return err
}
func (c *eventsController) processState(ctx tg.Context, usrState *domain.UserStateDTO) error {
	switch usrState.State {
	case defines.StateCreateCategoryWaitingName:
		return c.createCategoryWaitingName(ctx, usrState)
	case defines.StateCreateCategoryWaitingEmoji:
		return c.createCategoryWaitingEmoji(ctx, usrState)
	case defines.StateCreateWalletWaitingName:
		return c.createWalletWaitingName(ctx, usrState)
	case defines.StateCreateWalletWaitingAmount:
		return c.createWalletWaitingAmount(ctx, usrState)
	}

	return nil
}

// States handlers
func (c *eventsController) createCategoryWaitingName(ctx tg.Context, usrState *domain.UserStateDTO) error {
	userID := ctx.Sender().ID
	categoryName := ctx.Message().Text

	cat, ok := usrState.Data.(domain.APICategoryCreateRequest)
	if !ok {
		cat = domain.APICategoryCreateRequest{}
	}
	cat.Name = categoryName
	usrState.Data = cat

	// Update state
	usrState.State = defines.StateCreateCategoryWaitingEmoji

	err := c.usrStateSvc.UpdateByUserID(userID, usrState)
	if err != nil {
		return err
	}

	// Respond to the user
	err = ctx.Send(
		defines.MessageCreateCategoryWaitingEmoji,
		tg.ModeMarkdown,
	)

	return err
}
func (c *eventsController) createCategoryWaitingEmoji(ctx tg.Context, usrState *domain.UserStateDTO) error {
	userID := ctx.Sender().ID
	categoryEmoji := ctx.Message().Text

	var cat domain.APICategoryCreateRequest
	err := maptostruct.Convert(usrState.Data, &cat)
	if err != nil {
		return err
	}

	cat.Emoji = categoryEmoji

	_, err = c.catSvc.Create(userID, cat.Name, cat.Emoji)
	if err != nil {
		return err
	}

	err = c.usrStateSvc.DeleteByUserID(userID)
	if err != nil {
		return err
	}

	// Respond to the user
	err = ctx.Send(
		fmt.Sprintf(defines.MessageCreateCategorySuccess, cat.Name, cat.Emoji),
		tg.ModeMarkdown,
	)

	return err
}
func (c *eventsController) createWalletWaitingName(ctx tg.Context, usrState *domain.UserStateDTO) error {
	userID := ctx.Sender().ID
	walletName := ctx.Message().Text

	w, ok := usrState.Data.(domain.APIWalletCreateRequest)
	if !ok {
		w = domain.APIWalletCreateRequest{}
	}
	w.Name = walletName
	usrState.Data = w

	// Update state
	usrState.State = defines.StateCreateWalletWaitingAmount

	err := c.usrStateSvc.UpdateByUserID(userID, usrState)
	if err != nil {
		return err
	}

	// Respond to the user
	err = ctx.Send(
		defines.MessageCreateWalletWaitingAmount,
		tg.ModeMarkdown,
	)

	return err
}
func (c *eventsController) createWalletWaitingAmount(ctx tg.Context, usrState *domain.UserStateDTO) error {
	userID := ctx.Sender().ID
	balanceStr := ctx.Message().Text
	balanceStr = strings.Replace(balanceStr, "$", "", 1)
	balanceStr = strings.Replace(balanceStr, ",", ".", 1)
	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return err
	}

	var w domain.APIWalletCreateRequest
	err = maptostruct.Convert(usrState.Data, &w)
	if err != nil {
		return err
	}

	createdAt := time.Unix(ctx.Message().Unixtime, 0)

	w.InitialAmount = balance

	_, err = c.walletsSvc.Create(userID, w.Name, w.InitialAmount, &createdAt)
	if err != nil {
		return err
	}

	err = c.usrStateSvc.DeleteByUserID(userID)
	if err != nil {
		return err
	}

	// Respond to the user
	err = ctx.Send(
		fmt.Sprintf(defines.MessageCreateWalletSuccess, w.Name),
		tg.ModeMarkdown,
	)

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

func (c *eventsController) beginCreateTransaction(ctx tg.Context) error {
	userID := ctx.Sender().ID

	amount, description, err := c.getParametersFromMessage(ctx.Message().Text)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Create and set status to next step: wallet selection
	err = c.usrStateSvc.Create(userID, amount, description, ctx.Message().Time(), int64(ctx.Message().ID), defines.StateCreateTransactionSelectingWallet)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	token, err := c.usrSvc.GetToken(userID)
	wallets, err := c.walletsSvc.GetAll(token)

	kb := buildWalletsKeyboard(wallets)

	// Respond to the user
	err = ctx.Send(
		"¿Con qué billetera?", // TODO Mover a defines
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
func (c *eventsController) walletSelection(ctx tg.Context, txnStatus *domain.UserStateDTO) error {
	userID := ctx.Sender().ID
	token, err := c.usrSvc.GetToken(userID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	categories, err := c.catSvc.GetAll(token)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Update and change status to next step: category selection
	walletID := strings.Replace(ctx.Callback().Data, "\f", "", 1)

	var txn domain.APITransactionCreateRequest
	err = maptostruct.Convert(txnStatus.Data, &txn)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}
	txn.WalletID = walletID
	txnStatus.Data = txn
	txnStatus.State = defines.StateCreateTransactionSelectingCategory
	err = c.usrStateSvc.UpdateByUserID(userID, txnStatus)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	kb := buildCategoriesKeyboard(categories)

	// Respond to the user
	err = ctx.EditOrSend(
		"¿De cuál categoría?", // TODO Mover a defines
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
func (c *eventsController) categorySelection(ctx tg.Context, txnStatus *domain.UserStateDTO) error {
	userID := ctx.Sender().ID
	token, err := c.usrSvc.GetToken(userID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	var txn domain.APITransactionCreateRequest
	err = maptostruct.Convert(txnStatus.Data, &txn)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	catID := strings.Replace(ctx.Callback().Data, "\f", "", 1)

	cat, err := c.catSvc.GetByID(catID, token)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	wal, err := c.walletsSvc.GetByID(txn.WalletID, token)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Create transaction
	err = c.txnSvc.AddTransaction(
		txn.MsgID,
		txn.Amount,
		txn.Description,
		wal.ID,
		cat.ID,
		txn.CreatedAt,
		token,
	)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	// Delete status
	err = c.usrStateSvc.DeleteByUserID(userID)
	if err != nil {
		c.errorHandler(ctx, err)
		return err
	}

	var msg string

	if txn.Amount > 0 {
		msg = fmt.Sprintf(defines.MessageAddPaymentResponse,
			txn.Description,
			cat.Name,
			formatFloat(txn.Amount),
			wal.Name,
		)
	} else {
		msg = fmt.Sprintf(defines.MessageAddMoneyResponse,
			txn.Description,
			cat.Name,
			formatFloat(txn.Amount),
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

func buildWalletsKeyboard(wallets *[]domain.APIWalletGetResponse) *tg.ReplyMarkup {
	kb := &tg.ReplyMarkup{}

	var btns []tg.Btn

	for _, w := range *wallets {
		text := w.Name
		if w.Emoji != "" {
			text = w.Emoji + " " + text
		}
		btns = append(btns, kb.Data(text, w.ID))
	}
	rows := kb.Split(2, btns)
	kb.Inline(rows...)

	return kb
}
func buildCategoriesKeyboard(c *[]domain.APICategoryGetResponse) *tg.ReplyMarkup {
	kb := &tg.ReplyMarkup{}

	var btns []tg.Btn

	for _, c := range *c {
		text := c.Name
		if c.Emoji != "" {
			text = c.Emoji + " " + text
		}
		btns = append(btns, kb.Data(text, c.ID))
	}
	rows := kb.Split(2, btns)
	kb.Inline(rows...)

	return kb
}
