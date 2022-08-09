package bot

import (
	"log"
	"loquegasto-telegram/internal/controller"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/service"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	tgbot "gopkg.in/tucnak/telebot.v2"
)

var bot *tgbot.Bot

func New() *tgbot.Bot {
	var err error
	bot, err = tgbot.NewBot(tgbot.Settings{
		Token: os.Getenv(defines.EnvTelegramToken),
		Poller: &tgbot.LongPoller{
			Timeout: 30 * time.Second,
		},
		Verbose: true,
	})
	if err != nil {
		log.Fatalln(err)
	}

	mapCommands()

	return bot
}

func mapCommands() {
	// Init rest client
	rc := resty.New()

	// Init repositories
	txnRepo := repository.NewTransactionsRepository(rc)
	usersRepo := repository.NewUsersRepository(rc)
	walletsRepo := repository.NewWalletsRepository(rc)

	// Init services
	txnSrv := service.NewTransactionsService(txnRepo)
	usersSrv := service.NewUsersService(usersRepo)
	walletsSrv := service.NewWalletsService(walletsRepo)
	sheetsSrv := service.NewSheetsService()

	// Init controllers
	cmdCtrl := controller.NewCommandsController(bot, txnSrv, usersSrv, walletsSrv)
	parserCtrl := controller.NewParserController(bot, txnSrv, walletsSrv, sheetsSrv)
	grpCtrl := controller.NewGroupsController(bot)

	// Commands
	bot.Handle(defines.CommandStart, cmdCtrl.Start)
	bot.Handle(defines.CommandHelp, cmdCtrl.Help)
	bot.Handle(defines.CommandGetWallets, cmdCtrl.GetWallets)
	bot.Handle(defines.CommandCreateWallet, cmdCtrl.CreateWallet)
	bot.Handle(defines.CommandPing, cmdCtrl.Ping)
	bot.Handle(defines.CommandAddTransaction, cmdCtrl.AddTransaction)

	// Parser
	bot.Handle(tgbot.OnText, parserCtrl.Parse)
	bot.Handle(tgbot.OnEdited, parserCtrl.ParseEdited)

	// Group events
	bot.Handle(tgbot.OnAddedToGroup, grpCtrl.Start)
	bot.Handle(tgbot.OnUserJoined, grpCtrl.RegisterUsers)

}
