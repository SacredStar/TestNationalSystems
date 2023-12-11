package loging

import (
	"github.com/rs/zerolog"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	*zerolog.Logger
}

var instance Logger
var once sync.Once

func StartLog(level zerolog.Level, pathToLogFile string) *Logger {
	once.Do(func() {
		logFile, err := os.OpenFile(pathToLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		//Using multi-writer to write Stdout(customize time format) and log file
		mw := io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC850}, logFile)
		// Customize field name
		zerolog.TimestampFieldName = "t"
		zerolog.LevelFieldName = "l"
		zerolog.MessageFieldName = "m"
		zerolog.CallerFieldName = "c"

		innerLogger := zerolog.New(mw).With().Timestamp().Caller().Logger()
		innerLogger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(level)
		log.SetOutput(innerLogger)
		instance = Logger{
			Logger: &innerLogger,
		}
	})
	return &instance
}
