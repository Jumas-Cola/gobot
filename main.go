package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Jumas-Cola/gonzofilter"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "+", "\\+")
	s = strings.ReplaceAll(s, "-", "\\-")
	s = strings.ReplaceAll(s, "#", "\\#")
	return s
}

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
			// user = c.Sender()
			text = c.Text()
			msg  = c.Message()
		)

		if len([]rune(text)) < 30 {
			return nil
		}

		res := gonzofilter.ClassifyMessage(text, "hamspam.db")
		if res == "SPAM" {
			b.Delete(msg)
			return c.Send(fmt.Sprintf("Mabe spam: ||%s||", escapeString(text)),
				&tele.SendOptions{
					ParseMode: tele.ModeMarkdownV2,
				})
		}

		return nil
	})

	b.Start()
}
