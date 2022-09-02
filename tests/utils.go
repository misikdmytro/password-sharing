package tests

import (
	"context"

	"github.com/misikdmitriy/password-sharing/database"
	"github.com/misikdmitriy/password-sharing/model"
)

func MigrateDatabase(c context.Context, f database.DbFactory) error {
	db, close, err := f.InitDB(c)
	if err != nil {
		return err
	}
	defer close()

	err = db.Migrator().DropTable(&model.Password{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.Password{})
	if err != nil {
		return err
	}

	return nil
}
