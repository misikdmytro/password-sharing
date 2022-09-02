package health

import (
	"context"
	"errors"

	"github.com/misikdmitriy/password-sharing/database"
	"github.com/misikdmitriy/password-sharing/logger"
	"go.uber.org/zap"
)

type pgHealthCheck struct {
	factory       database.DbFactory
	loggerFactory logger.LoggerFactory
}

func NewPgHealthCheck(factory database.DbFactory, loggerFactory logger.LoggerFactory) HealthCheck {
	return &pgHealthCheck{
		factory:       factory,
		loggerFactory: loggerFactory,
	}
}

func (pg *pgHealthCheck) Check(c context.Context) (bool, error) {
	db, dbClose, err := pg.factory.InitDB(c)
	if err != nil {
		return false, err
	}
	defer dbClose()

	if err := db.Exec("SELECT 1").Error; err != nil {
		appLogger, loggerClose, err := pg.loggerFactory.NewLogger()
		if err != nil {
			return false, err
		}
		defer loggerClose()

		appLogger.Error("error on pg health check",
			zap.Error(err))

		return false, errors.New("pg health check failed")
	}

	return true, nil
}
