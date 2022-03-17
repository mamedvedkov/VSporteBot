package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) sendErrLog(update tgbotapi.Update, err error) {
	txt := fmt.Sprintf("Name: %s %s %s\n",
		update.SentFrom().FirstName, update.SentFrom().LastName, update.SentFrom().UserName)
	txt += fmt.Sprintf("Msg: %s\n", update.Message.Text)
	txt += fmt.Sprintf("Err: %v", err)

	_, err = b.api.Send(tgbotapi.NewMessage(int64(b.maintainerId), "`"+txt+"`"))
	if err != nil {
		log.Printf("cant send log to maintainer, err: %v , log: %s", err, txt)
	}
}
