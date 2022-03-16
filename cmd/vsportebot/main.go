package main

import (
	"context"
	"github.com/mamedvedkov/VSporteBot/internal/bot"
	"github.com/mamedvedkov/tools/env"
)

const (
	envToken        = "TELEGRAM_BOT_TOKEN"
	envMaintainerId = "MAINTAINER_ID"
)

func main() {
	token := env.Get(envToken).MustString()
	maintainerId := env.Get(envMaintainerId).MustInt()

	_bot := bot.New(token, maintainerId, 5, newRegisterStore())

	err := _bot.HandleMessages(context.Background())
	panic(err)
}
