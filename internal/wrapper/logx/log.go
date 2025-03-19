package logx

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	log     *LogX
	logOnce sync.Once
)

type LogX struct {
	*logrus.Logger
}

func Provide() *LogX {
	logOnce.Do(func() {
		// ============== global log ==============

		/// ENV: LOG_LEVEL
		// "panic"
		// "fatal"
		// "error"
		// "warn", "warning"
		// "info"
		// "debug"
		// "trace"
		logLevel, ok := os.LookupEnv("LOG_LEVEL")

		// LOG_LEVEL not set, let's default to debug
		if !ok {
			logLevel = "debug"
		}

		logrus.SetLevel(logrus.DebugLevel)

		// parse string, this is built-in feature of logrus
		logrusLevel, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrusLevel = logrus.DebugLevel
		}

		// set global log level
		logrus.SetLevel(logrusLevel)

		// TimestampFormat: "2006-01-02 150405"
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		logrus.SetOutput(os.Stdout)

		// ============== old log ==============

		logrusLogger := logrus.New()
		logrusLogger.SetLevel(logrus.DebugLevel)

		// parse string, this is built-in feature of logrus
		// logrusLevel, err := logrus.ParseLevel(logLevel)
		// if err != nil {
		// 	logrusLevel = logrus.DebugLevel
		// }

		// set global log level
		logrusLogger.SetLevel(logrusLevel)

		// TimestampFormat: "2006-01-02 150405"
		logrusLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		logrusLogger.SetOutput(os.Stdout)

		log = &LogX{logrusLogger}
	})

	return log
}
