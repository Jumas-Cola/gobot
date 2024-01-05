package bayesAntispam

import (
	"log/slog"

	gonzofilter "github.com/Jumas-Cola/gonzofilter"
	"gopkg.in/telebot.v3"
)

type BayesAntispam struct {
	databasePath string
}

func GetExtension() BayesAntispam {
	return BayesAntispam{databasePath: "../hamspam.db"}
}

func (ba BayesAntispam) RegisterHandlers(b *telebot.Bot) []telebot.Command {
    // TODO: Добавить команды для включения/выключения, 
    // сохранять настройки для каждого чата в MongoDB
	cmds := []telebot.Command{}

	b.Handle(telebot.OnText, func(c telebot.Context) error {
		var (
			text = c.Text()
			msg  = c.Message()
		)

		if len([]rune(text)) < 20 {
			return nil
		}

		res := gonzofilter.ClassifyMessage(text, ba.databasePath)
		if res == "SPAM" {
			b.Delete(msg)
			slog.Warn("Spam detected: " + text)
		}

		return nil
	})

	return cmds
}
