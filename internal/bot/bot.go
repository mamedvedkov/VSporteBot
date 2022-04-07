package bot

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/go-logr/logr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RegistredUsersStore interface {
	IsRegistred(ctx context.Context, id string) (name string, ok bool)
}

type PaymentsStore interface {
	CurrentMonthPayment(ctx context.Context, name string) (string, error)
	PastMonthPayment(ctx context.Context, name string) (string, error)
	PaymentDetail(ctx context.Context, name string) (string, error)
}

type RequisitesStore interface {
	RequisitesInfo(ctx context.Context, name string) (string, error)
}

type Bot struct {
	logger logr.Logger

	maintainerId int64
	finDepId     int64

	api  *tgbotapi.BotAPI
	pool *sync.Pool

	updates tgbotapi.UpdatesChannel

	userStore       RegistredUsersStore
	paymentsStore   PaymentsStore
	requisitesStore RequisitesStore
}

func New(
	logger logr.Logger,
	token string,
	maintainerId int64,
	finDepId int64,
	poolsize int,
	userStore RegistredUsersStore,
	paymentsStore PaymentsStore,
	requisitesStore RequisitesStore,
) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	// api.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := api.GetUpdatesChan(u)

	bot := &Bot{
		logger:       logger.WithName("tgBot"),
		maintainerId: maintainerId,
		finDepId:     finDepId,
		api:          api,
		pool: func() *sync.Pool {
			var pool sync.Pool

			for poolsize > 0 {
				pool.Put(struct{}{})
				poolsize--
			}

			return &pool
		}(),
		updates:         updates,
		userStore:       userStore,
		paymentsStore:   paymentsStore,
		requisitesStore: requisitesStore,
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
	if update.CallbackQuery != nil {
		b.handleCallback(ctx, update.CallbackQuery)
		return
	}

	if update.Message == nil {
		return
	}

	ok, err := b.authorize(update.Message.From.ID)
	if err != nil {
		b.logger.Error(err, "auth error")
		return
	}

	if !ok {
		b.handleStart(ctx, update.Message.From.ID, update.Message.Chat.ID, false)
		return
	}

	if b.isPaymentDocumentUrl(update.Message.Text) {
		b.forwardToFinDep(update)
	}

	if update.Message.Text == "Главное меню" {
		b.handleMainByText(ctx, update)
	}

	if update.Message.IsCommand() {
		b.handleCommand(ctx, update)
		return
	}
}

func (b *Bot) authorize(id int64) (bool, error) {
	_, ok := b.userStore.IsRegistred(context.Background(), strconv.Itoa(int(id)))
	return ok, nil
}

func (b *Bot) handleCommand(ctx context.Context, update tgbotapi.Update) {
	ok, err := b.authorize(update.Message.From.ID)
	if err != nil {
		b.logger.Error(err, "auth error")
		return
	}

	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	if !ok {
		b.handleStart(ctx, userId, chatId, false)
		return
	}

	switch update.Message.Command() {
	case "start":
		b.handleStart(ctx, userId, chatId, true)
	}
}

func (b *Bot) handleMainByText(ctx context.Context, update tgbotapi.Update) {
	_delete := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	b.api.Send(_delete)

	b.handleStart(ctx, update.Message.From.ID, update.Message.Chat.ID, true)
}

func (b *Bot) handleStart(ctx context.Context, userId, chatId int64, authorized bool) {
	if !authorized {
		err := b.handleFirstStart(userId, chatId)
		if err != nil {
			b.logger.Error(err, "send error")
		}

		return
	}

	msg := tgbotapi.NewMessage(chatId,
		"Привет! Это бот финансового департамента команды VSporte."+
			" Пожалуйста, выбери интересующее тебя меню")
	msg.ReplyMarkup = mainMenuIK()
	_, err := b.api.Send(msg)
	if err != nil {
		b.logger.Error(err, "send error")
	}
}

func (b *Bot) handleFirstStart(userId, chatId int64) error {
	msg := tgbotapi.NewMessage(
		chatId,
		"Сообщите ваш ID руководству производственного департамента: Кириллу Разумовскому или Ольге Скоробогатовой"+
			" для регистрации в системе\n\nПосле получения подтверждения о регистрации введите /start еще раз\n\n"+
			"Нажмите для копирования, ваш ID: "+
			fmt.Sprintf("`%v`", userId),
	)
	msg.ReplyMarkup = testKeyboard
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := b.api.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

var documentUrl = regexp.
	MustCompile("^(https://lknpd.nalog.ru\\/api\\/v1\\/receipt\\/)+(\\d{12})+(/)+(\\w{10})+\\/print+$")

func (b *Bot) isPaymentDocumentUrl(url string) bool {
	return documentUrl.MatchString(url)
}

func (b *Bot) forwardToFinDep(update tgbotapi.Update) {
	fwd := tgbotapi.NewForward(b.finDepId, update.Message.From.ID, update.Message.MessageID)

	_, err := b.api.Send(fwd)
	if err != nil {
		b.sendErrLog(update, err)

		return
	}

	del := tgbotapi.NewDeleteMessage(update.Message.From.ID, update.Message.MessageID)
	b.api.Send(del)

	msg := tgbotapi.NewMessage(update.Message.From.ID, "Ваша ссылка передана в финансовый департамент")
	_, err = b.api.Send(msg)
	if err != nil {
		b.sendErrLog(update, err)

		return
	}
}
