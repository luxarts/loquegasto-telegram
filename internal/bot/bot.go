package bot

import (
	"github.com/go-redis/redis/v9"
	"log"
	"loquegasto-telegram/internal/controller"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/repository"
	"loquegasto-telegram/internal/service"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	tgbot "gopkg.in/telebot.v3"
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
	restClient := resty.New()
	// Init redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv(defines.EnvRedisHost),
	})

	// Init repositories
	txnRepo := repository.NewTransactionsRepository(restClient)
	usersRepo := repository.NewUsersRepository(restClient)
	walletsRepo := repository.NewWalletsRepository(restClient)
	usrStateRepo := repository.NewUserStateRepository(redisClient)
	catRepo := repository.NewCategoriesRepository(restClient)

	// Init services
	txnSvc := service.NewTransactionsService(txnRepo)
	usersSvc := service.NewUsersService(usersRepo)
	walletsSvc := service.NewWalletsService(walletsRepo)
	usrStateSvc := service.NewUserStateService(usrStateRepo)
	catSvc := service.NewCategoriesService(catRepo)

	// Init controllers
	cmdCtrl := controller.NewCommandsController(bot, txnSvc, usersSvc, walletsSvc, usrStateSvc)
	evtCtrl := controller.NewEventsController(bot, txnSvc, usrStateSvc, walletsSvc, catSvc)

	// Commands
	bot.Handle(defines.CommandStart, cmdCtrl.Start)
	bot.Handle(defines.CommandHelp, cmdCtrl.Help)
	bot.Handle(defines.CommandGetWallets, cmdCtrl.GetWallets)
	bot.Handle(defines.CommandCreateWallet, cmdCtrl.CreateWallet)
	bot.Handle(defines.CommandCreateCategory, cmdCtrl.CreateCategory)
	bot.Handle(defines.CommandCancel, cmdCtrl.Cancel)
	bot.Handle(defines.CommandPing, cmdCtrl.Ping)

	// Events
	bot.Handle(tgbot.OnText, evtCtrl.Parse)
	//bot.Handle(tgbot.OnEdited, evtCtrl.ParseEdited)
	bot.Handle(tgbot.OnCallback, evtCtrl.Process)
}
