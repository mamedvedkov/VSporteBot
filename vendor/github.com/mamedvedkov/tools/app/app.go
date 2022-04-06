package app

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/mamedvedkov/tools/processes"
	"github.com/rotisserie/eris"
	"go.uber.org/zap"
)

type App struct {
	ctx context.Context

	workers []processes.Process
	l       logr.Logger
	closers []func()
}

type Option func(app *App)

func WithContext(ctx context.Context) Option {
	return func(app *App) {
		app.ctx = ctx
	}
}

func WithCustomLogger(logger logr.Logger) Option {
	return func(app *App) {
		app.l = logger
	}
}

func NewAppWithOpts(opts ...Option) *App {
	app := NewApp()

	for idx := range opts {
		opts[idx](app)
	}

	return app
}

func NewApp() *App {
	zapLogger, err := zap.NewProduction(zap.AddCaller())
	if err != nil {
		panic(err)
	}

	a := &App{
		ctx: context.Background(),
		l:   zapr.NewLogger(zapLogger),
	}

	a.AddWorkers(SignalNotify)

	return a
}

func (a *App) Logger() logr.Logger {
	return a.l
}

func (a *App) AddClosers(cls ...func()) *App {
	a.closers = append(cls, a.closers...)

	return a
}

func (a *App) Close() {
	for _, c := range a.closers {
		c()
	}
}

func (a *App) AddWorkers(rs ...processes.Process) *App {
	a.workers = append(a.workers, rs...)

	return a
}

func (a *App) Run() {
	a.l.Info("app run")

	Exit(func() error {
		return processes.RunParallelAndWait(a.ctx, a.workers...)
	}, a.l)
}

// Exit функция в зависимости от результат errFn выполнить os.Exit.
//
// Если errFn вернет ошибку не типа ErrorOSSignal, то os.Exit(1), иначе return.
// nolint:interfacer // дублирование интерфейса logr.
func Exit(errFn func() error, l logr.Logger) {
	err := errFn()

	if err == nil {
		return
	}

	var es *ErrorOSSignal

	if isSignal := errors.As(err, &es); isSignal {
		l.Info("app is stopped by signal", "signal", es.signal.String())

		return
	}

	l.Error(err, "app is stopped by error")

	panic(err)
}

func RecoverExit(l logr.Logger) {
	if r := recover(); r != nil {
		l.WithValues("stacktrace",
			eris.ToJSON(eris.Errorf(""), true)["root"].(map[string]interface{})["stack"]).
			Error(fmt.Errorf("%s", r), "recovered panic")
		os.Exit(1)
	}
}
