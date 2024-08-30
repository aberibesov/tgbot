package commands

import (
	"encoding/json"
	"github.com/aberibesov/tgbot/internal/service/product"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

//var registeredCommands = map[string]func(c *Commander, msg *tgbotapi.Message){}

type Commander struct {
	Bot            *tgbotapi.BotAPI
	ProductService *product.Service
	waitNew        bool
	updateIdx      int
	confirmDelete  bool
}

func (c *Commander) resetFlags(exclude string) {
	if exclude != "waitNew" {
		c.waitNew = false
	}
	if exclude != "confirmDelete" {
		c.confirmDelete = false
	}
	if exclude != "updateIdx" {
		c.updateIdx = -1
	}
}

type commandData struct {
	Command string `json:"command"`
	Idx     int    `json:"idx"`
}

func NewCommander(bot *tgbotapi.BotAPI, productService *product.Service) *Commander {
	return &Commander{
		Bot:            bot,
		ProductService: productService,
		updateIdx:      -1,
	}
}

func (c *Commander) HandleUpdate(update tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recovered from panic %+v", err)
		}
	}()

	if update.CallbackQuery != nil {
		parsedData := commandData{}
		json.Unmarshal([]byte(update.CallbackQuery.Data), &parsedData)

		switch parsedData.Command {
		case "new":
			c.New(update.CallbackQuery.Message)
		case "edit":
			c.Edit(update.CallbackQuery.Message, parsedData.Idx)
		case "confirm":
			c.confirmDelete = true
			c.Delete(update.CallbackQuery.Message, parsedData.Idx)
		case "delete":
			c.Delete(update.CallbackQuery.Message, parsedData.Idx)
		case "cancel":
			c.resetFlags("")
			c.List(update.CallbackQuery.Message)
		default:
			c.resetFlags("")
			c.Default(update.CallbackQuery.Message)
		}
		return
	}

	if update.Message == nil {
		return
	}

	/*command, ok := registeredCommands[update.Message.Command()]
	if ok {
		command(c, update.Message)
	} else {
		c.Default(update.Message)
	}
	*/
	if update.Message != nil { // If we got a message
		switch update.Message.Command() {
		case "help":
			c.waitNew = false
			c.Help(update.Message)
		case "list":
			c.waitNew = false
			c.List(update.Message)
		case "get":
			c.waitNew = false
			c.Get(update.Message)
		case "new":
			c.New(update.Message)
		case "delete":
			c.waitNew = false
			c.Delete(update.Message, 0)
		case "edit":
			c.waitNew = false
			c.Edit(update.Message, 0)
		default:
			if c.waitNew {
				c.New(update.Message)
			} else if c.updateIdx >= 0 {
				c.Edit(update.Message, c.updateIdx)
			} else {
				c.Default(update.Message)
			}
		}
	}
}
