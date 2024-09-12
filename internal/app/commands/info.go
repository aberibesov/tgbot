package commands

import (
	"fmt"
	"github.com/aberibesov/tgbot/internal/service/product"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func (c *Commander) Info(inputMessage *tgbotapi.Message, settings *UserSettings) {
	if settings.Step == 1 {
		c.handleWaitingForLS(inputMessage, settings)
	} else {
		args := inputMessage.CommandArguments()

		ls, err := strconv.Atoi(args)

		if ls > 0 && err == nil {
			c.sendRequestByLs(inputMessage, ls, settings)
			return
		}
		settings.Step = 1
		c.promptForClientAccount(inputMessage)
	}
}

func (c *Commander) handleWaitingForLS(inputMessage *tgbotapi.Message, settings *UserSettings) {
	ls, err := strconv.Atoi(inputMessage.Text)
	if err != nil {
		c.sendMessage(inputMessage.Chat.ID, "Лицевой счет должен содержать только цифры")
		return
	}

	c.sendRequestByLs(inputMessage, ls, settings)
}

func (c *Commander) sendRequestByLs(inputMessage *tgbotapi.Message, ls int, settings *UserSettings) {
	response, err := c.ProductService.GetInfo(ls)
	if err != nil {
		log.Println(err)
		c.sendMessage(inputMessage.Chat.ID, "Internal error occurred.")
		return
	}

	settings.Step = 0
	settings.State = Default
	msgText := c.formatClientInfo(response)
	c.sendMessage(inputMessage.Chat.ID, msgText)
}

func (c *Commander) promptForClientAccount(inputMessage *tgbotapi.Message) {
	c.sendMessage(inputMessage.Chat.ID, "Введите лицевой счет клиента:")
}

func (c *Commander) formatClientInfo(response *product.Response) string {
	if response.Code == "OK" {
		return fmt.Sprintf(
			"Фамилия: %s\nИмя: %s\nОтчество: %s\nАдрес: %s\nЛогин: %s\nПароль: %s\nДата регистрации: %s\nБаланс: %.1f\nТариф: %s\n",
			response.Client.Surname,
			response.Client.Name,
			response.Client.Patronymic,
			response.Client.Address,
			response.Client.Login,
			response.Client.Password,
			response.Client.DateReg,
			response.Client.Money,
			response.Client.Tariff,
		)
	}
	return response.Code
}

func (c *Commander) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := c.Bot.Send(msg); err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}
