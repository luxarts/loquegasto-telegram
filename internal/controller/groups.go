package controller

import (
	"fmt"
	tg "gopkg.in/tucnak/telebot.v2"
	"log"
	"loquegasto-telegram/internal/defines"
)

type GroupsController interface {
	Start(m *tg.Message)
	RegisterUsers(m *tg.Message)
}
type groupsController struct {
	bot *tg.Bot
}

func NewGroupsController(bot *tg.Bot) GroupsController {
	return &groupsController{bot: bot}
}

func (c *groupsController) Start(m *tg.Message) {
	c.botRespond(m, defines.MessageStartGroup)
}
func (c *groupsController) RegisterUsers(m *tg.Message) {
	if len(m.UsersJoined) == 1 {
		c.bot.Send(m.Chat, "Hola @"+m.UsersJoined[0].Username+"!")
	} else {
		c.bot.Send(m.Chat, "Hola a todos!")
	}
}

// Utils
func (c *groupsController) errorHandler(m *tg.Message, err error) {
	log.Println(err)
	_, err = c.bot.Send(m.Chat, defines.MessageError, tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *groupsController) errorHandlerResponse(m *tg.Message, err error) {
	log.Println(err)
	_, err = c.bot.Send(m.Chat, fmt.Sprintf(defines.MessageErrorResponse, err.Error()), tg.ModeMarkdown)
	if err != nil {
		log.Println(err)
	}
}
func (c *groupsController) botRespond(m *tg.Message, msg string) {
	if _, err := c.bot.Send(m.Chat, msg, tg.ModeMarkdown); err != nil {
		c.errorHandler(m, err)
	}
}
