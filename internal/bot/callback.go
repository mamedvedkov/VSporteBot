package bot

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	callbackData := callback.Data

	var edit tgbotapi.EditMessageTextConfig

	switch callbackData {
	case infoMsg:
		edit = b.handleInfo(ctx, callback)
	case mainMenu:
		edit = b.handleMain(ctx, callback)
	case myPayment:
		edit = b.handleMyPayment(ctx, callback)
	case faq:
		edit = b.handleFaq(ctx, callback)
	case personalData:
		edit = b.handlePersonal(ctx, callback)
	case detailPyments:
		edit = b.handleDetailPayments(ctx, callback)
	case monthPayments + currentMonth:
		edit = b.handleMonth(ctx, callback, currentMonth)
	case monthPayments + prevMonth:
		edit = b.handleMonth(ctx, callback, prevMonth)
	}

	iK, ok := callBackRoutes[callbackData]
	if !ok {
		log.Printf("no route for callback: %v\n", callbackData)
	}
	edit.ReplyMarkup = &iK

	_, err := b.api.Send(edit)
	if err != nil {
		log.Printf("err callback send: %v\n", err)
	}
}

func (b *Bot) handleInfo(ctx context.Context, callBack *tgbotapi.CallbackQuery) tgbotapi.EditMessageTextConfig {
	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		"ебало опусти",
	)

	return edit
}

func (b *Bot) handleMyPayment(ctx context.Context, callBack *tgbotapi.CallbackQuery) tgbotapi.EditMessageTextConfig {
	var text string

	text = "Затычка под поле с оплатой"

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		text,
	)

	return edit
}

func (b *Bot) handleFaq(ctx context.Context, callBack *tgbotapi.CallbackQuery) tgbotapi.EditMessageTextConfig {
	var text string

	text = "faq"

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		text,
	)

	return edit
}

func (b *Bot) handlePersonal(ctx context.Context, callBack *tgbotapi.CallbackQuery) tgbotapi.EditMessageTextConfig {
	var text string

	role, _ := b.userStore.IsRegistred(ctx, int(callBack.From.ID))

	text = fmt.Sprintf("Ваш айди и роль:\n %v, %s",
		callBack.From.ID, role)

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		text,
	)

	return edit
}

func (b *Bot) handleDetailPayments(ctx context.Context, callBack *tgbotapi.CallbackQuery) tgbotapi.EditMessageTextConfig {
	var text string

	text = "ссылка на детализацию"

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		text,
	)

	return edit
}

func (b *Bot) handleMonth(ctx context.Context, callBack *tgbotapi.CallbackQuery, month string) tgbotapi.EditMessageTextConfig {
	var text string

	if month == currentMonth {
		text = "текущий месяц"
	}

	if month == prevMonth {
		text = "прошлый месяц"
	}

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		text,
	)

	return edit
}

func (b *Bot) handleMain(ctx context.Context, callBack *tgbotapi.CallbackQuery) tgbotapi.EditMessageTextConfig {
	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		"Выберете меню",
	)

	return edit
}
