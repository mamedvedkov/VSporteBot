package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-logr/logr"
)

type Env struct {
	key     string
	value   string
	def     interface{}
	exists  bool
	needLog bool
	logger  logr.Logger
}

func Get(key string) *Env {
	return GetWithOption(key)
}

type Option func(env *Env) *Env

func GetWithOption(key string, opts ...Option) *Env {
	value, ok := os.LookupEnv(key)

	env := &Env{
		key:    key,
		value:  value,
		exists: ok,
	}

	for i := range opts {
		env = opts[i](env)
	}

	return env
}

func WithLogging(logger logr.Logger) Option {
	return func(env *Env) *Env {
		env.needLog = true
		env.logger = logger

		return env
	}
}

const (
	envFoundMsg      = "ENV %s FOUND = %s"
	envNotFoundMsg   = "ENV %s NOTFOUND, DEFAULT = %v"
	noValuePanicText = "no value for key=%s"
)

func logIfNeeded(env *Env) {
	if !env.needLog || !env.logger.Enabled() {
		return
	}

	msg := envNotFoundMsg
	val := env.def
	if env.exists {
		msg = envFoundMsg
		val = env.value
	}

	env.logger.Info(fmt.Sprintf(msg, env.key, val))
}

func panicIfNotExist(env *Env) {
	if !env.exists {
		panic(fmt.Sprintf(noValuePanicText, env.key))
	}
}

func (e *Env) String(def string) string {
	e.def = def
	logIfNeeded(e)

	if e.exists {
		return e.value
	}

	return def
}

func (e *Env) MustString() string {
	panicIfNotExist(e)

	return e.String("")
}

func (e *Env) Int(def int) int {
	e.def = def
	logIfNeeded(e)

	if e.exists {
		val, err := strconv.Atoi(e.value)
		if err != nil {
			panic(err)
		}
		return val
	}

	return def
}

func (e *Env) MustInt() int {
	panicIfNotExist(e)

	return e.Int(0)
}

// TODO
//func (e *Env) Float(def float64) float64 {
//
//}
//
// TODO
//func (e *Env) Time(def time.Duration) time.Duration {
//
//}
