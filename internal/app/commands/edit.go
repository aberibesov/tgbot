package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (c *Commander) Edit(inputMessage *tgbotapi.Message, settings *UserSettings, idx int) {
	msgText := ""

	product, err := c.ProductService.Get(idx)
	if err != nil {
		log.Printf("Fail to get product with %d: %v", idx, err)
		return
	}

	if settings.Idx >= 0 && settings.Step == 1 {
		msgText += product.Title + " изменен на " + inputMessage.Text
		c.ProductService.Update(settings.Idx, inputMessage.Text)
		settings.Idx = -1
		settings.Step = 0
		settings.State = Default
	} else {
		msgText += "Введите новое значение для " + product.Title
		settings.Idx = idx
		settings.Step = 1
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID,
		msgText,
	)

	_, errSend := c.Bot.Send(msg)
	if errSend != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}
