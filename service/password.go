package service

import (
	"github.com/jackc/pgconn"
	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/database"
	"github.com/misikdmitriy/password-sharing/helper"
	"github.com/misikdmitriy/password-sharing/model"
)

type PasswordService interface {
	CreateLinkFromPassword(pwd string) (string, error)
}

type passwordService struct {
	dbFactory database.DbFactory
	conf      *config.Config
	rf        helper.RandomGeneratorFactory
}

func NewPasswordService(dbFactory database.DbFactory, conf *config.Config,
	rf helper.RandomGeneratorFactory) PasswordService {
	return &passwordService{
		dbFactory: dbFactory,
		conf:      conf,
		rf:        rf,
	}
}

const pgUniqueViolationCode = "23505"

func (s *passwordService) CreateLinkFromPassword(pwd string) (string, error) {
	db, err := s.dbFactory.InitDB()
	if err != nil {
		return "", err
	}

	conn, err := db.DB()
	defer conn.Close()

	for {
		rg := s.rf.NewRandomGenerator()
		link, err := rg.RandomString(s.conf.App.LinkLength)

		command := db.Save(model.NewPassword(link, pwd))
		if err := command.Error; err != nil {
			pgErr, ok := err.(*pgconn.PgError)
			if ok && pgErr.Code == pgUniqueViolationCode {
				continue
			}

			return "", err
		}

		return link, err
	}
}
