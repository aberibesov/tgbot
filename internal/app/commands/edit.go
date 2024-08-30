package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (c *Commander) Edit(inputMessage *tgbotapi.Message, idx int) {
	msgText := ""

	product, err := c.ProductService.Get(idx)
	if err != nil {
		log.Printf("Fail to get product with %d: %v", idx, err)
		return
	}

	if c.updateIdx >= 0 {
		msgText += product.Title + " изменен на " + inputMessage.Text
		c.ProductService.Update(c.updateIdx, inputMessage.Text)
		c.updateIdx = -1
	} else {
		msgText += "Введите новое значение для " + product.Title
		c.updateIdx = idx
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID,
		msgText,
	)

	_, errSend := c.Bot.Send(msg)
	if errSend != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}
