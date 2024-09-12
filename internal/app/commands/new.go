package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (c *Commander) New(inputMessage *tgbotapi.Message, settings *UserSettings) {
	msgText := ""
	if settings.Step == 1 {
		msgText += "Добавлен новый элемент: " + inputMessage.Text
		c.ProductService.Add(inputMessage.Text)
		settings.Step = 0
		settings.State = Default
	} else {
		settings.Step = 1
		msgText += "Введите название продукта:"
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID,
		msgText,
	)

	_, err := c.Bot.Send(msg)
	if err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}
