package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Jumas-Cola/gonzofilter"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env")
	}

	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tele.OnText, func(c tele.Context) error {
		var (
			text = c.Text()
			msg  = c.Message()
		)

		if len([]rune(text)) < 20 {
			return nil
		}

		res := gonzofilter.ClassifyMessage(text, "hamspam.db")
		if res == "SPAM" {
			b.Delete(msg)
			slog.Warn("Spam detected: " + text)
		}

		return nil
	})

	b.Start()
}
