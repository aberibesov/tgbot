package commands

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func (c *Commander) Get(inputMessage *tgbotapi.Message) {
	args := inputMessage.CommandArguments()

	idx, err := strconv.Atoi(args)
	if err != nil {
		log.Println("Invalid argument", args)
		return
	}

	product, err := c.ProductService.Get(idx)
	if err != nil {
		log.Printf("Fail to get product with %d: %v", idx, err)
		return
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID,
		fmt.Sprintf(product.Title),
	)

	serializedComUpd, _ := json.Marshal(commandData{"edit", idx})
	serializedComDel, _ := json.Marshal(commandData{"delete", idx})

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Edit", string(serializedComUpd)),
			tgbotapi.NewInlineKeyboardButtonData("Delete", string(serializedComDel)),
		),
	)

	c.Bot.Send(msg)
}
