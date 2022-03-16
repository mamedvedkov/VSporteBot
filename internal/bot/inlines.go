package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var exampleNumericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

var startNK = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonSwitch("switch", "switch"),
		tgbotapi.NewInlineKeyboardButtonData("data", "data"),
		tgbotapi.NewInlineKeyboardButtonURL("url", "http://1.com"),
	),
)
