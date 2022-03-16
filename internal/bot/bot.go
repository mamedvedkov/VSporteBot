package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"sync"
)

type RegistredUsersStore interface {
	IsRegistred(ctx context.Context, id int) (role string, ok bool)
	Register(ctx context.Context, id int, role string)
}

type Bot struct {
	maintainerId int

	api  *tgbotapi.BotAPI
	pool *sync.Pool

	updates tgbotapi.UpdatesChannel

	userStore RegistredUsersStore
}

func New(
	token string,
	maintainerId int,
	poolsize int,
	userStore RegistredUsersStore,
) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	api.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := api.GetUpdatesChan(u)

	bot := &Bot{
		maintainerId: maintainerId,
		api:          api,
		pool: func() *sync.Pool {
			var pool sync.Pool

			for poolsize > 0 {
				pool.Put(struct{}{})
				poolsize--
			}

			return &pool

		}(),
		updates:   updates,
		userStore: userStore,
	}

	return bot
}

func (b *Bot) HandleMessages(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-b.updates:
			b.pool.Get()
			go func(u tgbotapi.Update) {
				b.handleUpdate(ctx, u)
				b.pool.Put(struct{}{})
			}(update)
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message.IsCommand() {
		b.handleCommand(ctx, update)
	}
}

func (b *Bot) handleCommand(ctx context.Context, update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		b.handleStart(ctx, update)
	case "register":
		b.handleRegistration(ctx, update)
	}
}

func (b *Bot) handleStart(ctx context.Context, update tgbotapi.Update) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	_, ok := b.userStore.IsRegistred(ctx, int(userId))
	if !ok {
		err := b.handleFirstStart(userId, chatId)
		if err != nil {
			b.sendErrLog(update, err)
		}

		return
	}

}

func (b *Bot) handleFirstStart(userId, chatId int64) error {
	msg := tgbotapi.NewMessage(
		chatId,
		"Сообщите ваш айди менеджеру для добавления в белый список\n\nНажмите для копирования: "+
			fmt.Sprintf("`%v`", userId),
	)

	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := b.api.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleRegistration(ctx context.Context, update tgbotapi.Update) {
	id, err := strconv.Atoi(update.Message.CommandArguments())
	if err != nil {
		b.sendErrLog(update, err)
	}

	b.userStore.Register(ctx, id, "mock-role")
}
