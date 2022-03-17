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

const (
	mainMenu      = "mainMenu"
	myPayment     = "myPayments"
	faq           = "faq"
	personalData  = "personalData"
	infoMsg       = "info"
	detailPyments = "detailPyments"
	monthPayments = "monthPayments"
)

const (
	currentMonth = "cur"
	prevMonth    = "prev"
)

var callBackRoutes = map[string]tgbotapi.InlineKeyboardMarkup{
	mainMenu:                     mainMenuIK,
	faq:                          onlyBackIK,
	personalData:                 onlyBackIK,
	infoMsg:                      onlyBackIK,
	myPayment:                    myPaymentsIK,
	monthPayments + currentMonth: myMonthlyPaymentsIK,
	monthPayments + prevMonth:    myMonthlyPaymentsIK,
	detailPyments:                detailIK,
}

var mainMenuIK = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Моя зп", myPayment),
		tgbotapi.NewInlineKeyboardButtonData("У меня тупой вопрос", faq),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Мои данные", personalData),
		tgbotapi.NewInlineKeyboardButtonData("Как это работает", infoMsg),
	),
)

var onlyBackIK = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
	),
)

var myPaymentsIK = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", monthPayments+currentMonth),
		tgbotapi.NewInlineKeyboardButtonData("Прошлый месяц", monthPayments+prevMonth),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Детализация", detailPyments),
		tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
	),
)

var myMonthlyPaymentsIK = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", myPayment),
		tgbotapi.NewInlineKeyboardButtonData("Детализация", detailPyments),
	),
)

var detailIK = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", myPayment),
		tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
	),
)
