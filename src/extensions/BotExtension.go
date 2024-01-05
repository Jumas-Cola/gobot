package extensions

import "gopkg.in/telebot.v3"

type BotExtension interface {
	RegisterHandlers(*telebot.Bot) []telebot.Command
}
