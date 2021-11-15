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
	tg "gopkg.in/tucnak/telebot.v2"
)

var bot *tg.Bot

func New() *tg.Bot {
	var err error
	bot, err = tg.NewBot(tg.Settings{
		Token: os.Getenv(defines.EnvTelegramToken),
		Poller: &tg.LongPoller{
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

	// Init controllers
	cmdCtrl := controller.NewCommandsController(bot, txnSrv, usersSrv, walletsSrv)
	parserCtrl := controller.NewParserController(bot, txnSrv, walletsSrv)

	// Commands
	bot.Handle(defines.CommandStart, cmdCtrl.Start)
	bot.Handle(defines.CommandHelp, cmdCtrl.Help)
	bot.Handle(defines.CommandWallets, cmdCtrl.Wallets)
	bot.Handle(defines.CommandCreateWallet, cmdCtrl.CreateWallet)
	bot.Handle(defines.CommandPing, cmdCtrl.Ping)

	// Parser
	bot.Handle(tg.OnText, parserCtrl.Parse)
	bot.Handle(tg.OnEdited, parserCtrl.ParseEdited)
}
