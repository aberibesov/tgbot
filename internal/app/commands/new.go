package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (c *Commander) New(inputMessage *tgbotapi.Message, settings *UserSettings) {
	msgText := ""
	if settings.WaitNew {
		msgText += "Добавлен новый элемент: " + inputMessage.Text
		c.ProductService.Add(inputMessage.Text)
		settings.WaitNew = false
	} else {
		settings.WaitNew = true
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
