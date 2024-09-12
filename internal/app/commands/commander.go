package commands

import (
	"encoding/json"
	"github.com/aberibesov/tgbot/internal/service/product"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

//var registeredCommands = map[string]func(c *Commander, msg *tgbotapi.Message){}

type State int

const (
	New State = iota
	List
	Edit
	Delete
	Info
	Get
	Default
)

type UserSettings struct {
	State State
	Step  int
	Idx   int
}

type Commander struct {
	Bot            *tgbotapi.BotAPI
	ProductService *product.Service
	MapUsers       map[int64]UserSettings
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
			userSettings.State = New
			c.New(update.CallbackQuery.Message, &userSettings)
		case "edit":
			userSettings.State = Edit
			c.Edit(update.CallbackQuery.Message, &userSettings, parsedData.Idx)
		case "confirm":
			userSettings.State = Delete
			userSettings.Step = 1
			c.Delete(update.CallbackQuery.Message, &userSettings, parsedData.Idx)
		case "delete":
			userSettings.State = Delete
			c.Delete(update.CallbackQuery.Message, &userSettings, parsedData.Idx)
		case "cancel":
			userSettings.State = Delete
			c.List(update.CallbackQuery.Message)
		default:
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
		c.MapUsers[update.Message.From.ID] = UserSettings{Idx: -1}
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
		userSettings.State = Default
		c.Help(update.Message)
	case "list":
		userSettings.State = List
		c.List(update.Message)
	case "get":
		userSettings.State = Get
		c.Get(update.Message)
	case "new":
		userSettings.State = New
		c.New(update.Message, &userSettings)
	case "delete":
		userSettings.State = Delete
		c.Delete(update.Message, &userSettings, -1)
	case "edit":
		userSettings.State = Edit
		c.Edit(update.Message, &userSettings, -1)
	case "info":
		userSettings.State = Info
		c.Info(update.Message, &userSettings)
	default:
		switch userSettings.State {
		case New:
			c.New(update.Message, &userSettings)
		case Edit:
			c.Edit(update.Message, &userSettings, userSettings.Idx)
		case Info:
			c.Info(update.Message, &userSettings)
		default:
			c.Default(update.Message)
		}
	}
	c.MapUsers[update.Message.From.ID] = userSettings
}
