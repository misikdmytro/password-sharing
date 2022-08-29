package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/database"
	"github.com/misikdmitriy/password-sharing/helper"
	"github.com/misikdmitriy/password-sharing/logger"
	"github.com/misikdmitriy/password-sharing/tests"
)

func TestCreateLinkFromPasswordShouldDoIt(t *testing.T) {
	c := config.CreateEmpty()
	c.Database.ConnectionString = "inmemdb"
	c.Database.Provider = "sqlite"

	c.App.LinkLength = 8

	log, err := logger.NewLogger()
	if err != nil {
		t.Error(err)
	}

	dbf := database.NewFactory(c, log)
	err = tests.MigrateDatabase(dbf)
	if err != nil {
		t.Error(err)
	}

	rf := helper.NewRandomFactory()
	s := NewPasswordService(dbf, c, rf, log)

	result, err := s.CreateLinkFromPassword(uuid.New().String())
	if err != nil {
		t.Error(err)
	}

	if len(result) != c.App.LinkLength {
		t.Errorf("expected password length to be %d but was %d", c.App.LinkLength, len(result))
	}
}
