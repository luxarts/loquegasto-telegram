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

	// Init services
	txnSrv := service.NewTransactionsService(txnRepo)

	// Init controllers
	cmdCtrl := controller.NewCommandsController(bot, txnSrv)
	parserCtrl := controller.NewParserController(bot, txnSrv)

	// Commands
	bot.Handle(defines.CommandStart, cmdCtrl.Start)
	bot.Handle(defines.CommandHelp, cmdCtrl.Help)
	bot.Handle(defines.CommandPing, cmdCtrl.Ping)

	// Parser
	bot.Handle(tg.OnText, parserCtrl.Parse)
	// TODO Agregar handler tg.OnEdit que actualice la información del pago
}
