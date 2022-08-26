package tests

import (
	"github.com/misikdmitriy/password-sharing/database"
	"github.com/misikdmitriy/password-sharing/model"
)

func MigrateDatabase(f database.DbFactory) error {
	db, err := f.InitDB()
	if err != nil {
		return err
	}

	conn, err := db.DB()
	if err != nil {
		return err
	}

	defer conn.Close()

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
