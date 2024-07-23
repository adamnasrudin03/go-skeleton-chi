package driver

import (
	"github.com/adamnasrudin03/go-skeleton-chi/app/configs"

	"github.com/sirupsen/logrus"
)

func Logger(config *configs.Configs) *logrus.Logger {
	logger := logrus.New()

	level := logrus.InfoLevel
	switch config.App.Env {
	case "dev":
		level = logrus.TraceLevel
	}

	logger.SetLevel(level)

	return logger
}
