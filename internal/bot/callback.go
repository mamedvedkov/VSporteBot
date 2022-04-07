package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mamedvedkov/VSporteBot/internal"
)

func isFaq(callbackData string) bool {
	return strings.Contains(callbackData, "faq")
}

func (b *Bot) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	defer func() {
		b.api.Send(tgbotapi.NewCallback(callback.ID, ""))
	}()

	callbackData := callback.Data

	var edit *tgbotapi.EditMessageTextConfig
	var err error

	if isFaq(callbackData) {
		edit, err = b.handleFaq(ctx, callback)
	}

	switch callbackData {
	case infoMsg:
		edit = b.handleInfo(ctx, callback)
	case mainMenu:
		edit = b.handleMain(ctx, callback)
	case myPayment:
		edit = b.handleMyPayment(ctx, callback)
	case personalData:
		edit, err = b.handlePersonal(ctx, callback)
	case changePersonalData:
		edit = b.handleChangePersonal(ctx, callback)
	case detailPyments:
		edit, err = b.handleDetailPayments(ctx, callback)
	case currentMonthPayments:
		edit, err = b.handleCurrentMonth(ctx, callback)
	case pastMonthPayments:
		edit, err = b.handlePastMonth(ctx, callback)
	case uploadFile:
		edit = b.handleUploadFile(callback)
	case compensation:
		edit = b.handleCompensation(callback)
	case compensationRules:
		b.handleCompensationRules(callback)
		return
	case szInstruction:
		b.handleSzInstruction(callback)
		return
	}

	if err != nil {
		b.logger.Error(err, "cant get edit")
		return
	}

	if edit == nil {
		b.logger.Info("no edit")
		return
	}

	name, ok := b.userStore.IsRegistred(ctx, strconv.Itoa(int(callback.From.ID)))
	if !ok {
		b.logger.Info("no user found")
	}

	detailUrl, err := b.paymentsStore.PaymentDetail(ctx, name)
	if err != nil {
		b.logger.Error(err, "cant get detail: %w")
	}

	inlineBuilder, ok := callBackRoutes(detailUrl)[callbackData]
	if !ok {
		log.Printf("no route for callback: %v\n", callbackData)
		return
	}

	iK := inlineBuilder()
	edit.ReplyMarkup = &iK

	_, err = b.api.Send(edit)
	if err != nil {
		log.Printf("err callback send: %v\n", err)
	}
}

func (b *Bot) handleInfo(ctx context.Context, callBack *tgbotapi.CallbackQuery) *tgbotapi.EditMessageTextConfig {
	var text string

	text = "Еще раз привет! Это бот команды VSporte, созданный с целью помочь тебе в коммуникации с финансовым отделом." +
		" Наша команда постаралась учесть все самые частые вопросы и автоматизировать многие процессы," +
		" чтобы дать тебе возможность оперативно и 24/7 получать любую интересующую тебя информацию.\n\n" +
		"Мы всегда рады обратной связи!"

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		text,
	)

	return &edit
}

func (b *Bot) handleMyPayment(ctx context.Context, callBack *tgbotapi.CallbackQuery) *tgbotapi.EditMessageTextConfig {
	var text string

	text = "В этом разделе ты можешь посмотреть данные за предыдущий месяц," +
		" *предварительные* данные за текущий и получить ссылку на детализированную ведомость"

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		text,
	)

	edit.ParseMode = tgbotapi.ModeMarkdown

	return &edit
}

func (b *Bot) handlePersonal(ctx context.Context, callBack *tgbotapi.CallbackQuery) (*tgbotapi.EditMessageTextConfig, error) {
	var text string

	name, _ := b.userStore.IsRegistred(ctx, strconv.Itoa(int(callBack.From.ID)))

	text = fmt.Sprintf("Ваш айди:\t%v\nФИО:\t%s\n",
		callBack.From.ID, name)

	requisites, err := b.requisitesStore.RequisitesInfo(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("cant get requisites: %w", err)
	}

	text += requisites

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		fmt.Sprintf("`%s`", text),
	)

	edit.ParseMode = tgbotapi.ModeMarkdown

	return &edit, nil
}

func (b *Bot) handleChangePersonal(ctx context.Context, callBack *tgbotapi.CallbackQuery) *tgbotapi.EditMessageTextConfig {
	text := "Если у тебя поменялись данные, пожалуйста, направь актуальные реквизиты личным сообщением в финансовый отдел" +
		" @VSporte_Finance\n\nВнесение изменений займет от 1 до 2 рабочих дней"

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		fmt.Sprintf(text))

	return &edit
}

func (b *Bot) handleDetailPayments(ctx context.Context, callBack *tgbotapi.CallbackQuery) (*tgbotapi.EditMessageTextConfig, error) {
	name, ok := b.userStore.IsRegistred(ctx, strconv.Itoa(int(callBack.From.ID)))
	if !ok {
		b.logger.Info("no user found")
	}

	var text string

	text, err := b.paymentsStore.PaymentDetail(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("cant get detail: %w", err)
	}

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		text,
	)

	return &edit, nil
}

func (b *Bot) handlePastMonth(ctx context.Context, callBack *tgbotapi.CallbackQuery) (*tgbotapi.EditMessageTextConfig, error) {
	name, ok := b.userStore.IsRegistred(ctx, strconv.Itoa(int(callBack.From.ID)))
	if !ok {
		b.logger.Info("no user found")
	}

	paymentSum, err := b.paymentsStore.PastMonthPayment(ctx, name)
	if err != nil {
		b.logger.Error(err, "cant get current moth payments:")

		paymentSum = "0 руб"
	}

	text := "За %s ты заработал: %s.\n\n" +
		"Если есть вопросы - посмотри детализацию по ссылке ниже или обратись к своему руководителю.\n"

	pastMonthName, err := internal.GetPastRussianMonth(int(time.Now().Month()))
	if err != nil {
		return nil, fmt.Errorf("month error: %w", err)
	}

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		fmt.Sprintf(text, pastMonthName, paymentSum),
	)

	return &edit, nil
}

func (b *Bot) handleCurrentMonth(ctx context.Context, callBack *tgbotapi.CallbackQuery) (*tgbotapi.EditMessageTextConfig, error) {
	name, ok := b.userStore.IsRegistred(ctx, strconv.Itoa(int(callBack.From.ID)))
	if !ok {
		b.logger.Info("no user found")
	}

	paymentSum, err := b.paymentsStore.CurrentMonthPayment(ctx, name)
	if err != nil {
		b.logger.Error(err, "cant get current moth payments:")

		paymentSum = "0 руб"
	}

	text := "*Расчёт не окончательный*,первичные данные собраны автоматически. Окончательный расчёт за %s" +
		" будет сформирован после 5-го числа следующего месяца\n*Предварительно*: %s."

	currentMonthName, err := internal.GetRussianMonth(int(time.Now().Month()))
	if err != nil {
		return nil, fmt.Errorf("month error: %w", err)
	}

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		fmt.Sprintf(text, currentMonthName, paymentSum),
	)

	edit.ParseMode = tgbotapi.ModeMarkdown

	return &edit, nil
}

func (b *Bot) handleMain(ctx context.Context, callBack *tgbotapi.CallbackQuery) *tgbotapi.EditMessageTextConfig {
	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		"Привет! Это бот финансового департамента команды VSporte. Пожалуйста, выбери интересующее тебя меню",
	)

	return &edit
}

func (b *Bot) handleUploadFile(callBack *tgbotapi.CallbackQuery) *tgbotapi.EditMessageTextConfig {
	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		"Для отправки чека в фин отдел отправьте ссылку на чек в этот чат",
	)

	return &edit
}

func (b *Bot) handleCompensation(callback *tgbotapi.CallbackQuery) *tgbotapi.EditMessageTextConfig {
	text := "Если ты потратил свои личные деньги на нужды компании или проекта ты можешь подать на компенсацию этих" +
		" расходов.\n\nВыплаты происходят 3 раза в неделю после подтверждения руководства." +
		" Подробнее о компенсациях в регламенте"

	edit := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		text,
	)

	return &edit
}

func (b *Bot) handleCompensationRules(callback *tgbotapi.CallbackQuery) {
	doc := tgbotapi.NewInputMediaDocument(tgbotapi.FilePath("./compensation.pdf"))
	media := tgbotapi.NewMediaGroup(callback.From.ID, []interface{}{doc})

	b.api.Send(media)

	b.handleStart(context.Background(), callback.From.ID, callback.Message.Chat.ID, true)
}

func (b *Bot) handleSzInstruction(callback *tgbotapi.CallbackQuery) {
	doc := tgbotapi.NewInputMediaDocument(tgbotapi.FilePath("./szInstruction.pdf"))
	media := tgbotapi.NewMediaGroup(callback.From.ID, []interface{}{doc})

	b.api.Send(media)

	b.handleStart(context.Background(), callback.From.ID, callback.Message.Chat.ID, true)
}

func (b *Bot) handleFaq(ctx context.Context, callBack *tgbotapi.CallbackQuery) (*tgbotapi.EditMessageTextConfig, error) {
	var edit *tgbotapi.EditMessageTextConfig
	var err error

	switch callBack.Data {
	case faq:
		edit = b.handleFaqMain(callBack)
	case faq1:
		edit, err = b.handleFaq1(callBack)
	case faq2:
		edit, err = b.handleFaq2(callBack)
	case faq3:
		edit, err = b.handleFaq3(callBack)
	case faq4:
		edit = b.handleFaq4(callBack)
	case faq5:
		edit = b.handleFaq5(callBack)
	case faq51:
		edit = b.handleFaq51(callBack)
	case faq6:
		edit = b.handleFaq6(callBack)
	default:
		err = fmt.Errorf("no valid faq")
	}

	if err != nil {
		return nil, err
	}

	return edit, err
}

func (b *Bot) handleFaqMain(callback *tgbotapi.CallbackQuery) *tgbotapi.EditMessageTextConfig {
	var text string

	text = "Согласно статистике, Финансовому отделу в 88% случаев задают такие вопросы.\n\nВыбери свой"

	edit := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		text,
	)

	return &edit
}

func (b *Bot) handleFaq1(callback *tgbotapi.CallbackQuery) (*tgbotapi.EditMessageTextConfig, error) {
	var text string
	text = "Выплата гонораров за работу на проектах текущего месяца происходит с 10 по 20 число следующего месяца." +
		"\n\nНапример, проект был назначен на %s. Ты отработал свою смену, и соответствующая" +
		" строчка появилась в твоем отчете. Тогда месяцем выплаты гонорара за него будет %s." +
		" Период выплат с с 10-го по 20-е (включительно)"

	currentMonthName, err := internal.GetRussianMonth(int(time.Now().Month()))
	if err != nil {
		return nil, fmt.Errorf("month error: %w", err)
	}

	pastMonthName, err := internal.GetPastRussianMonth(int(time.Now().Month()))
	if err != nil {
		return nil, fmt.Errorf("month error: %w", err)
	}

	edit := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		fmt.Sprintf(text, pastMonthName, currentMonthName),
	)

	return &edit, nil
}

func (b *Bot) handleFaq2(callback *tgbotapi.CallbackQuery) (*tgbotapi.EditMessageTextConfig, error) {
	var text string
	text = "Твой отчет за %s является промежуточным, данные в нем обновятся 5-го числа" +
		" следующего месяца.\n\nНекоторые новые проекты должны быть обработаны вручную менеджментом производственного" +
		" департамента, чтобы в твоих данных появились ставки за такие проекты. Это относится и к надбавкам/премиям" +
		" и штрафам за %s\n\nЭтот процесс будет автоматизирован в ближайшее время," +
		" наберись терпения!"

	currentMonthName, err := internal.GetRussianMonth(int(time.Now().Month()))
	if err != nil {
		return nil, fmt.Errorf("month error: %w", err)
	}

	edit := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		fmt.Sprintf(text, currentMonthName, currentMonthName),
	)

	return &edit, nil
}

func (b *Bot) handleFaq3(callback *tgbotapi.CallbackQuery) (*tgbotapi.EditMessageTextConfig, error) {
	var text string
	text = "Твой отчет за %s является промежуточным, данные в нем обновятся 5-го числа" +
		" следующего месяца.\n\nНадбавки, штрафы и премии к регулярным задачам формируются по итогам месяца." +
		"  Этот процесс также будет автоматизирован в ближайшее время, наберись терпения!\n\nЕсли твой отчет" +
		" за %s не полный, пожалуйста, напиши руководителю в производственный" +
		" департамент или менеджеру проекта"

	currentMonthName, err := internal.GetRussianMonth(int(time.Now().Month()))
	if err != nil {
		return nil, fmt.Errorf("month error: %w", err)
	}

	pastMonthName, err := internal.GetPastRussianMonth(int(time.Now().Month()))
	if err != nil {
		return nil, fmt.Errorf("month error: %w", err)
	}

	edit := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		fmt.Sprintf(text, currentMonthName, pastMonthName),
	)

	return &edit, nil
}

func (b *Bot) handleFaq4(callBack *tgbotapi.CallbackQuery) *tgbotapi.EditMessageTextConfig {
	var text string
	text = "Средства на командировочные расходы выплачиваются в случае твоей работы на проекте за" +
		" пределами домашнего региона.\n\nВыплата суточных назначается руководством проекта," +
		" эти выплаты строго регламентированы. С регламентом можно ознакомиться по ссылке ниже."

	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		text,
	)

	return &edit
}

func (b *Bot) handleFaq5(callBack *tgbotapi.CallbackQuery) *tgbotapi.EditMessageTextConfig {
	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		"Если ты зарегистрирован как самозанятый, то при выплате гонорара тебе придет уведомление о"+
			" необходимости сформировать чек с корректной суммой и ИНН юридического лица,"+
			" на которое необходимо выставить чек.\n\nОтправь этот чек прямо мне в чат,"+
			" я перешлю его в финансовый отдел",
	)

	return &edit
}

func (b *Bot) handleFaq51(callBack *tgbotapi.CallbackQuery) *tgbotapi.EditMessageTextConfig {
	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		"Если тебе уже пришло уведомление о просьбе сформировать чек, а деньги еще не поступили на счет то,"+
			" скорее всего, выплата еще обрабатывается банком или сегодня выходной и банковские платежи не проходят."+
			" Подожди 2 рабочих дня, если деньги так и не поступили напиши о проблеме в аккаунт @VSporte_Finance",
	)

	return &edit
}

func (b *Bot) handleFaq6(callBack *tgbotapi.CallbackQuery) *tgbotapi.EditMessageTextConfig {
	edit := tgbotapi.NewEditMessageText(
		callBack.Message.Chat.ID,
		callBack.Message.MessageID,
		"Если у тебя есть другой вопрос,"+
			" чтобы связаться с финансовым отделом напиши в их официальный аккаунт @VSporte_Finance",
	)

	return &edit
}
