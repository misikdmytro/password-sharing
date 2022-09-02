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
	c := &config.Config{}
	c.Database.ConnectionString = "inmemdb"
	c.Database.Provider = "sqlite"
	c.Encrypt.Secret = "123456789123456789012345"
	c.Encrypt.IV = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	ctxt := context.Background()
	c.App.LinkLength = 8

	loggerFactory := logger.NewTestLoggerFactory()
	encoder := helper.NewEncoder(c)
	dbf := database.NewFactory(c, loggerFactory)
	err := tests.MigrateDatabase(ctxt, dbf)
	if err != nil {
		t.Error(err)
	}

	rf := helper.NewRandomFactory()
	s := NewPasswordService(dbf, c, rf, loggerFactory, encoder)

	result, err := s.CreateLinkFromPassword(ctxt, uuid.New().String())
	if err != nil {
		t.Error(err)
	}

	if len(result) != c.App.LinkLength {
		t.Errorf("expected password length to be %d but was %d", c.App.LinkLength, len(result))
	}
}
