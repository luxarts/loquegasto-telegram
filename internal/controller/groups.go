package controller

import (
	"fmt"
	tg "gopkg.in/telebot.v3"
	"log"
	"loquegasto-telegram/internal/defines"
)

type GroupsController interface {
	Start(ctx tg.Context) error
	RegisterUsers(ctx tg.Context) error
}
type groupsController struct {
	bot *tg.Bot
}

func NewGroupsController(bot *tg.Bot) GroupsController {
	return &groupsController{bot: bot}
}

func (c *groupsController) Start(ctx tg.Context) error {
	c.botRespond(ctx, defines.MessageStartGroup)
	return nil
}
func (c *groupsController) RegisterUsers(ctx tg.Context) error {
	if len(ctx.Message().UsersJoined) == 1 {
		c.bot.Send(ctx.Recipient(), "Hola @"+ctx.Message().UsersJoined[0].Username+"!")
	} else {
		c.bot.Send(ctx.Recipient(), "Hola a todos!")
	}
	return nil
}

// Utils
func (c *groupsController) errorHandler(ctx tg.Context, err error) {
	log.Println(err)
	_, err = c.bot.Send(ctx.Recipient(), defines.MessageError, tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *groupsController) errorHandlerResponse(ctx tg.Context, err error) {
	log.Println(err)
	_, err = c.bot.Send(ctx.Recipient(), fmt.Sprintf(defines.MessageErrorResponse, err.Error()), tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *groupsController) botRespond(ctx tg.Context, msg string) {
	if _, err := c.bot.Send(ctx.Recipient(), msg, tg.ModeMarkdown); err != nil {
		c.errorHandler(ctx, err)
	}
}
