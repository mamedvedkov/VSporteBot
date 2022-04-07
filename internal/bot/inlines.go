package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	mainMenu      = "mainMenu"
	myPayment     = "myPayments"
	faq           = "faq"
	faq1          = "faq1"
	faq2          = "faq2"
	faq3          = "faq3"
	faq4          = "faq4"
	faq5          = "faq5"
	faq51         = "faq51"
	faq6          = "faq6"
	szInstruction = "szInstruction"

	personalData         = "personalData"
	changePersonalData   = "changePersonalData"
	infoMsg              = "info"
	detailPyments        = "detailPyments"
	currentMonthPayments = "currentMonthPayments"
	pastMonthPayments    = "pastMonthPayments"
	compensation         = "compensation"
	compensationRules    = "compensationRules"
	uploadFile           = "uploadFile"
)

type inlineMarkupBuilder func() tgbotapi.InlineKeyboardMarkup

var callBackRoutes = func(detailsUrl string) map[string]inlineMarkupBuilder {
	return map[string]inlineMarkupBuilder{
		mainMenu:             mainMenuIK,
		faq:                  faqIK,
		personalData:         personalDataIK,
		changePersonalData:   onlyBackIK,
		infoMsg:              onlyBackIK,
		myPayment:            myPaymentsIK(detailsUrl),
		currentMonthPayments: myCurrentMonthPaymentsIK(detailsUrl),
		pastMonthPayments:    myPastMonthPaymentsIK(detailsUrl),
		uploadFile:           onlyBackIK,
		compensation:         compensationIK,
		faq1:                 faqBackIK,
		faq2:                 faqBackIK,
		faq3:                 faqBackIK,
		faq4:                 faq4IK,
		faq5:                 faq5IK,
		faq51:                faqBackIK,
		faq6:                 faqBackIK,
	}
}

var mainMenuIK = func() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Моя зп", myPayment),
			tgbotapi.NewInlineKeyboardButtonData("Мои данные", personalData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Компенсация расходов", compensation),
			tgbotapi.NewInlineKeyboardButtonData("Как это работает", infoMsg),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Задать вопрос", faq),
		),
	)
}

var onlyBackIK = func() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
		),
	)
}

var personalDataIK = func() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить данные", changePersonalData),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
		),
	)
}

var faqIK = func() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Когда зарплата", faq1),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Нет ставки за проект", faq2),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отчет не полный", faq3),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Суточные", faq4),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Самозанятым", faq5),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Связаться с финотделом", faq6),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
		),
	)
}

var myPaymentsIK = func(detailUrl string) inlineMarkupBuilder {
	return func() tgbotapi.InlineKeyboardMarkup {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", currentMonthPayments),
				tgbotapi.NewInlineKeyboardButtonData("Прошлый месяц", pastMonthPayments),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Детализация", detailUrl),
				tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
			),
		)
	}
}

var myCurrentMonthPaymentsIK = func(detailUrl string) inlineMarkupBuilder {
	return func() tgbotapi.InlineKeyboardMarkup {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Прошлый месяц", pastMonthPayments),
				tgbotapi.NewInlineKeyboardButtonURL("Детализация", detailUrl),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", myPayment),
				tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
			),
		)
	}
}

var myPastMonthPaymentsIK = func(detailUrl string) inlineMarkupBuilder {
	return func() tgbotapi.InlineKeyboardMarkup {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", currentMonthPayments),
				tgbotapi.NewInlineKeyboardButtonURL("Детализация", detailUrl),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", myPayment),
				tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
			),
		)
	}
}

var compensationIK = func() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Компенсировать рассходы",
				"https://docs.google.com/forms/d/e/1FAIpQLSfznc49oGRoRfPUp8JA5wgL7nT_keQquv7TrCY_0gJPkEafcA/viewform"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Регламент выплат", compensationRules),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
		),
	)
}

var faqBackIK = func() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", faq),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
		),
	)
}

var faq4IK = func() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Регламент выплаты суточных", compensationRules),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", faq),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
		),
	)
}

var faq5IK = func() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Чек просят денег нет", faq51),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Как зарегестрироваться самозанятым", szInstruction),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", faq),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", mainMenu),
		),
	)
}

var testKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Главное меню"),
	),
)
