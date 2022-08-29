package service

import (
	"github.com/jackc/pgconn"
	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/database"
	"github.com/misikdmitriy/password-sharing/helper"
	"github.com/misikdmitriy/password-sharing/logger"
	"github.com/misikdmitriy/password-sharing/model"
)

type PasswordService interface {
	GetPasswordFromLink(string) (*string, error)
	CreateLinkFromPassword(string) (string, error)
}

type passwordService struct {
	dbFactory database.DbFactory
	conf      *config.Config
	rf        helper.RandomGeneratorFactory
	log       logger.Logger
}

func NewPasswordService(dbFactory database.DbFactory, conf *config.Config,
	rf helper.RandomGeneratorFactory,
	log logger.Logger) PasswordService {
	return &passwordService{
		dbFactory: dbFactory,
		conf:      conf,
		rf:        rf,
		log:       log,
	}
}

const pgUniqueViolationCode = "23505"

func (s *passwordService) CreateLinkFromPassword(pwd string) (string, error) {
	db, err := s.dbFactory.InitDB()
	if err != nil {
		s.log.Error("failed to init db")
		return "", err
	}

	conn, err := db.DB()
	defer conn.Close()

	for {
		rg := s.rf.NewRandomGenerator()
		link, err := rg.RandomString(s.conf.App.LinkLength)
		if err != nil {
			s.log.Error("error on randomizing",
				"length", s.conf.App.LinkLength,
				"error", err)

			return "", err
		}

		command := db.Save(model.NewPassword(link, pwd))
		if err := command.Error; err != nil {
			pgErr, ok := err.(*pgconn.PgError)
			if ok && pgErr.Code == pgUniqueViolationCode {
				s.log.Warn("retry after unique key violation")
				continue
			}

			s.log.Error("error on db command",
				"error", err)

			return "", err
		}

		s.log.Debug("link generated")
		return link, nil
	}
}

func (s *passwordService) GetPasswordFromLink(link string) (*string, error) {
	db, err := s.dbFactory.InitDB()
	if err != nil {
		s.log.Error("failed to init db")
		return nil, err
	}

	result := &model.Password{}
	query := db.Where(&model.Password{Link: link}).First(result)
	if err := query.Error; err != nil {
		s.log.Error("error on db query",
			"error", err)

		return nil, err
	}

	return &result.Password, nil
}
