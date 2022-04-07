package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) Send(id int64, msg string) {
	message := tgbotapi.NewMessage(id, msg)

	_, err := b.api.Send(message)
	if err != nil {
		b.logger.Error(err, "cant send message")
	}
}
