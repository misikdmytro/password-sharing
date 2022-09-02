package service

import (
	"context"
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

	ctxt := context.Background()
	c.App.LinkLength = 8

	log := logger.TestLogger()
	dbf := database.NewFactory(c, log)
	err := tests.MigrateDatabase(ctxt, dbf)
	if err != nil {
		t.Error(err)
	}

	rf := helper.NewRandomFactory()
	s := NewPasswordService(dbf, c, rf, log)

	result, err := s.CreateLinkFromPassword(ctxt, uuid.New().String())
	if err != nil {
		t.Error(err)
	}

	if len(result) != c.App.LinkLength {
		t.Errorf("expected password length to be %d but was %d", c.App.LinkLength, len(result))
	}
}
