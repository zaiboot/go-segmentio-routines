package logger

import (
	"os"
	"runtime"

	"github.com/rs/zerolog"
)

type LineInfoHook struct{}

func (h LineInfoHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	pc, _, _, ok := runtime.Caller(3)
	if ok {
		fn := runtime.FuncForPC(pc).Name()
		e.Str("FunctionName", fn)
	}
}

func InitLog(appName string, levelConfig string) zerolog.Logger {
	level, err := zerolog.ParseLevel(levelConfig)
	if err != nil {
		panic(err.Error())
	}

	log := zerolog.New(os.Stderr).Hook(LineInfoHook{}).With().
		Timestamp().
		Logger().Level(level)

	return log
}
