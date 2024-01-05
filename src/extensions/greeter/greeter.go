package greeter

import (
	"gopkg.in/telebot.v3"
)

type Greeter struct{}

func GetExtension() Greeter {
	return Greeter{}
}

func (g Greeter) RegisterHandlers(b *telebot.Bot) []telebot.Command {
	cmds := []telebot.Command{
		{Text: "hello", Description: "Say hello"},
	}

	b.Handle("/hello", g.handleHello)

	return cmds
}

func (g Greeter) handleHello(c telebot.Context) error {
	return c.Send("Hello, World!")
}
