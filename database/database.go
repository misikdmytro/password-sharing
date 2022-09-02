package database

import (
	"context"
	"fmt"

	"github.com/misikdmitriy/password-sharing/config"
	"go.uber.org/zap"
	"moul.io/zapgorm2"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbFactory interface {
	InitDB(context.Context) (*gorm.DB, func(), error)
}

type dbFactory struct {
	c   *config.Config
	log *zap.Logger
}

func NewFactory(conf *config.Config, log *zap.Logger) DbFactory {
	return &dbFactory{
		c:   conf,
		log: log,
	}
}

func (f *dbFactory) InitDB(c context.Context) (*gorm.DB, func(), error) {
	conn, err := f.createConnection()
	if err != nil {
		f.log.Error("cannot create db connection",
			zap.Error(err),
			zap.String("provider", f.c.Database.Provider),
		)

		return nil, nil, err
	}

	logger := zapgorm2.New(f.log)
	db, err := gorm.Open(*conn, &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		f.log.Error("cannot open gorm",
			zap.Error(err),
			zap.String("provider", f.c.Database.Provider),
		)

		return nil, nil, err
	}

	sql, err := db.DB()
	if err != nil {
		f.log.Error("failed on db get",
			zap.Error(err),
			zap.String("provider", f.c.Database.Provider),
		)

		return nil, nil, err
	}
	close := func() {
		sql.Close()
	}

	return db.WithContext(c), close, nil
}

func (f *dbFactory) createConnection() (*gorm.Dialector, error) {
	f.log.Debug("creating db connection",
		zap.String("provider", f.c.Database.Provider),
	)

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
