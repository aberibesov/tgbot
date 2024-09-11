package commands

import (
	"encoding/json"
	"github.com/aberibesov/tgbot/internal/service/product"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

//var registeredCommands = map[string]func(c *Commander, msg *tgbotapi.Message){}

type UserSettings struct {
	WaitNew       bool
	UpdateIdx     int
	ConfirmDelete bool
	WaitLS        bool
}

type Commander struct {
	Bot            *tgbotapi.BotAPI
	ProductService *product.Service
	MapUsers       map[int64]UserSettings
}

func (c *Commander) resetUserFlags(settings *UserSettings, exclude string) {
	if exclude != "waitNew" {
		settings.WaitNew = false
	}
	if exclude != "confirmDelete" {
		settings.ConfirmDelete = false
	}
	if exclude != "updateIdx" {
		settings.UpdateIdx = -1
	}
	if exclude != "waitLS" {
		settings.WaitLS = false
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
		MapUsers:       make(map[int64]UserSettings),
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
		userSettings := c.MapUsers[update.CallbackQuery.From.ID]
		switch parsedData.Command {
		case "new":
			c.New(update.CallbackQuery.Message, &userSettings)
		case "edit":
			c.Edit(update.CallbackQuery.Message, &userSettings, parsedData.Idx)
		case "confirm":
			userSettings.ConfirmDelete = true
			c.Delete(update.CallbackQuery.Message, &userSettings, parsedData.Idx)
		case "delete":
			c.Delete(update.CallbackQuery.Message, &userSettings, parsedData.Idx)
		case "cancel":
			c.resetUserFlags(&userSettings, "")
			c.List(update.CallbackQuery.Message)
		default:
			c.resetUserFlags(&userSettings, "")
			c.Default(update.CallbackQuery.Message)
		}
		c.MapUsers[update.CallbackQuery.From.ID] = userSettings
		return
	}

	if update.Message == nil {
		return
	}

	_, exists := c.MapUsers[update.Message.From.ID]
	if !exists {
		c.MapUsers[update.Message.From.ID] = UserSettings{UpdateIdx: -1}
	}

	/*command, ok := registeredCommands[update.Message.Command()]
	if ok {
		command(c, update.Message)
	} else {
		c.Default(update.Message)
	}
	*/
	// If we got a message
	userSettings := c.MapUsers[update.Message.From.ID]
	switch update.Message.Command() {
	case "help":
		userSettings.WaitNew = false
		c.Help(update.Message)
	case "list":
		userSettings.WaitNew = false
		c.List(update.Message)
	case "get":
		userSettings.WaitNew = false
		c.Get(update.Message)
	case "new":
		c.New(update.Message, &userSettings)
	case "delete":
		userSettings.WaitNew = false
		c.Delete(update.Message, &userSettings, -1)
	case "edit":
		userSettings.WaitNew = false
		c.Edit(update.Message, &userSettings, -1)
	case "info":
		c.Info(update.Message)
	default:
		if userSettings.WaitNew {
			c.New(update.Message, &userSettings)
		} else if userSettings.WaitLS {
			c.Info(update.Message)
		} else if userSettings.UpdateIdx >= 0 {
			c.Edit(update.Message, &userSettings, userSettings.UpdateIdx)
		} else {
			c.Default(update.Message)
		}
	}
	c.MapUsers[update.Message.From.ID] = userSettings
}
