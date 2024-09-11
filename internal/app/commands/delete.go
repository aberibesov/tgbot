package commands

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (c *Commander) Delete(inputMessage *tgbotapi.Message, settings *UserSettings, idx int) {
	msgText := ""

	product, err := c.ProductService.Get(idx)
	if err != nil {
		log.Printf("Fail to get product with %d: %v", idx, err)
		return
	}

	serializedComConfirm := []byte{}
	serializedComCancel := []byte{}

	if settings.ConfirmDelete {
		msgText += "Удален элемент:" + product.Title
		c.ProductService.Delete(idx)
		settings.ConfirmDelete = false
	} else {
		msgText += "Уверены, что хотите удалить элемент: " + product.Title
		serializedComConfirm, _ = json.Marshal(commandData{"confirm", idx})
		serializedComCancel, _ = json.Marshal(commandData{"cancel", idx})
		settings.ConfirmDelete = false
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID,
		msgText,
	)

	if len(serializedComCancel) > 0 {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Confirm", string(serializedComConfirm)),
				tgbotapi.NewInlineKeyboardButtonData("Cancel", string(serializedComCancel)),
			),
		)
	}

	_, errSend := c.Bot.Send(msg)
	if errSend != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}
