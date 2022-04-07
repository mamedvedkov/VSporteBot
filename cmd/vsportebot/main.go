package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/mamedvedkov/VSporteBot/internal/bot"
	"github.com/mamedvedkov/VSporteBot/internal/notifier"
	"github.com/mamedvedkov/VSporteBot/internal/repository/google_sheet"
	"github.com/mamedvedkov/tools/app"
	"github.com/mamedvedkov/tools/cron"
	"github.com/mamedvedkov/tools/env"
)

const (
	envToken        = "TELEGRAM_BOT_TOKEN"
	envMaintainerId = "MAINTAINER_ID"
	envFinDepId     = "FINDEP_ID"

	envPaymentSpreadsheetId             = "PAYMENT_SPREADSHEET_ID"
	envPaymentConsolidatedSpreadsheetId = "PAYMENT_CONSOLIDATED_SPREADSHEET_ID"
	envScheduleSpreadsheetId            = "SCHEDULE_SPREADSHEET_ID"
	envFinancialSpreadsheetId           = "FINANCIAL_SPREADSHEET_ID"

	envNotifySchedule     = "NOTIFY_SCHEDULE"
	defaultNotifySchedule = "0 30 10 */2 * *" // в 10:30 каждые 2 дня
)

const reloadInterval = 1 * time.Minute

func main() {
	_app := app.NewApp()
	defer _app.Close()

	_app.Logger().Info("read creds")
	fileData, err := os.ReadFile("./credentials.json")
	if err != nil {
		panic(err)
	}

	_app.Logger().Info("connect to google")

	ids := google_sheet.SpreadsheetsIds{
		Payment:             env.Get(envPaymentSpreadsheetId).MustString(),
		PaymentConsolidated: env.Get(envPaymentConsolidatedSpreadsheetId).MustString(),
		Schedule:            env.Get(envScheduleSpreadsheetId).MustString(),
		Financial:           env.Get(envFinancialSpreadsheetId).MustString(),
	}

	entityToInn := make(map[string]string)
	innData, err := os.ReadFile("./inn.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(innData, &entityToInn)
	if err != nil {
		panic(err)
	}

	googleClient := google_sheet.MustGoogle(_app.Logger(), fileData, ids, entityToInn)
	_app.AddWorkers(googleClient.ReloadSpreadSheets(reloadInterval))

	_app.Logger().Info("create bot")

	token := env.Get(envToken).MustString()
	maintainerId := env.Get(envMaintainerId).MustInt()
	findepId := env.Get(envFinDepId).MustInt()

	_bot := bot.New(
		_app.Logger(),
		token,
		int64(maintainerId),
		int64(findepId),
		5,
		googleClient,
		googleClient,
		googleClient,
	)

	_app.AddWorkers(_bot.HandleMessages)

	_app.Logger().Info("create cron")

	_cron := cron.New()

	_cron.MustAddJobs(
		cron.NewJob(env.Get(envNotifySchedule).String(defaultNotifySchedule), notifier.Notify(_bot, googleClient)),
	)
	_app.AddWorkers(_cron.Run)

	_app.Run()
}
