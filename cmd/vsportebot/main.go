package main

import (
	"github.com/mamedvedkov/VSporteBot/internal/bot"
	"github.com/mamedvedkov/tools/app"
	"github.com/mamedvedkov/tools/env"
)

const (
	envToken        = "TELEGRAM_BOT_TOKEN"
	envMaintainerId = "MAINTAINER_ID"
)

func main() {
	_app := app.NewApp()

	token := env.Get(envToken).MustString()
	maintainerId := env.Get(envMaintainerId).MustInt()

	_bot := bot.New(token, maintainerId, 5, newRegisterStore())

	_app.AddWorkers(_bot.HandleMessages)

	_app.Run()
}
