package database

import (
	"fmt"

	"github.com/misikdmitriy/password-sharing/config"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbFactory interface {
	InitDB() (*gorm.DB, error)
}

type dbFactory struct {
	c *config.Config
}

func NewFactory(conf *config.Config) DbFactory {
	return &dbFactory{
		c: conf,
	}
}

func (f *dbFactory) InitDB() (*gorm.DB, error) {
	conn, err := f.createConnection()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(*conn, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (f *dbFactory) createConnection() (*gorm.Dialector, error) {
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
