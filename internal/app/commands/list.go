package commands

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Commander) List(inputMessage *tgbotapi.Message) {
	outputMessageText := "Here all the Products: \n\n"

	products := c.ProductService.List()
	for _, p := range products {
		outputMessageText += p.Title + "\n"
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, outputMessageText)

	serializedComCr, _ := json.Marshal(commandData{"new", 0})

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("New", string(serializedComCr)),
		),
	)

	c.Bot.Send(msg)
}

/*func init() {
	registeredCommands["list"] = (*Commander).List
}*/
