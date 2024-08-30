package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (c *Commander) Help(inputMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID,
		"/help - help\n"+
			"/list - list Products\n"+
			"/get - get product by ID",
	)

	c.Bot.Send(msg)
}

/*func init() {
	registeredCommands["help"] = (*Commander).Help
}*/
