package health

import (
	"context"
	"errors"

	"github.com/misikdmitriy/password-sharing/database"
)

type pgHealthCheck struct {
	factory     database.DbFactory
	onUnhealthy func(err error)
}

func NewPgHealthCheck(factory database.DbFactory, onUnhealthy func(err error)) HealthCheck {
	return &pgHealthCheck{
		factory:     factory,
		onUnhealthy: onUnhealthy,
	}
}

func (pg *pgHealthCheck) Check(c context.Context) (bool, error) {
	db, err := pg.factory.InitDB(c)
	if err != nil {
		return false, err
	}

	if err := db.Exec("SELECT 1").Error; err != nil {
		pg.onUnhealthy(err)
		return false, errors.New("pg health check failed")
	}

	return true, nil
}
