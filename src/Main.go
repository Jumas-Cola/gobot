package main

import (
	"gobot/src/db"
	"gobot/src/extensions"
	"gobot/src/extensions/bayesAntispam"
	"gobot/src/extensions/greeter"
	"gobot/src/extensions/kikVote"
	"log/slog"

	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("Error loading .env")
	}

	pref := telebot.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	db.ConnectToMongoDB()

	// https://github.com/ti-bone/feedbackBot/blob/main/src/Main.go

	// Create extension instances
	extensions := []extensions.BotExtension{
		greeter.GetExtension(),
		bayesAntispam.GetExtension(),
		kikVote.GetExtension(),
	}

	var commands []telebot.Command

	for _, ext := range extensions {
		commands = append(commands, ext.RegisterHandlers(b)...)
	}

	err = b.SetCommands(commands)
	if err != nil {
		log.Println("Не удалось установить команды бота:", err)
	}

	b.Start()
}
