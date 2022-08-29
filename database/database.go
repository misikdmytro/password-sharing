package database

import (
	"fmt"

	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/logger"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbFactory interface {
	InitDB() (*gorm.DB, error)
}

type dbFactory struct {
	c   *config.Config
	log logger.Logger
}

func NewFactory(conf *config.Config, log logger.Logger) DbFactory {
	return &dbFactory{
		c:   conf,
		log: log,
	}
}

func (f *dbFactory) InitDB() (*gorm.DB, error) {
	conn, err := f.createConnection()
	if err != nil {
		f.log.Error("cannot create db connection",
			"provider", f.c.Database.Provider,
			"error", err)

		return nil, err
	}

	db, err := gorm.Open(*conn, &gorm.Config{})
	if err != nil {
		f.log.Error("cannot open gorm",
			"provider", f.c.Database.Provider,
			"error", err)

		return nil, err
	}

	return db, nil
}

func (f *dbFactory) createConnection() (*gorm.Dialector, error) {
	f.log.Debug("creating db connection",
		"provider", f.c.Database.Provider)

	switch f.c.Database.Provider {
	case "pg":
		conn := postgres.New(postgres.Config{
			DSN: f.c.Database.ConnectionString,
		})
		return &conn, nil
	case "sqlite":
		conn := sqlite.Open(f.c.Database.ConnectionString)
		return &conn, nil
	default:
		return nil, fmt.Errorf("cannot create %s connection", f.c.Database.Provider)
	}
}
