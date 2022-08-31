package service

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/database"
	pserror "github.com/misikdmitriy/password-sharing/error"
	"github.com/misikdmitriy/password-sharing/helper"
	"github.com/misikdmitriy/password-sharing/model"
	"go.uber.org/zap"
)

type PasswordService interface {
	GetPasswordFromLink(context.Context, string) (string, error)
	CreateLinkFromPassword(context.Context, string) (string, error)
}

type passwordService struct {
	dbFactory database.DbFactory
	conf      *config.Config
	rf        helper.RandomGeneratorFactory
	log       *zap.Logger
}

func NewPasswordService(dbFactory database.DbFactory, conf *config.Config,
	rf helper.RandomGeneratorFactory,
	log *zap.Logger) PasswordService {
	return &passwordService{
		dbFactory: dbFactory,
		conf:      conf,
		rf:        rf,
		log:       log,
	}
}

const pgUniqueViolationCode = "23505"

func (s *passwordService) CreateLinkFromPassword(c context.Context, pwd string) (string, error) {
	db, err := s.dbFactory.InitDB(c)
	if err != nil {
		return "", initDbError(s.log)
	}

	conn, err := db.DB()
	defer conn.Close()

	for {
		rg := s.rf.NewRandomGenerator()
		link, err := rg.RandomString(s.conf.App.LinkLength)
		if err != nil {
			const message = "error on randomizing"

			s.log.Error(message,
				zap.Error(err),
				zap.Int("length", s.conf.App.LinkLength),
			)

			return "", &pserror.PasswordSharingError{
				Code:    pserror.RandomizerError,
				Message: message,
			}
		}

		command := db.Save(model.NewPassword(link, pwd))
		if err := command.Error; err != nil {
			pgErr, ok := err.(*pgconn.PgError)
			if ok && pgErr.Code == pgUniqueViolationCode {
				s.log.Warn("retry after unique key violation")
				continue
			}

			const message = "error on db command"

			s.log.Error(message,
				zap.Error(err),
			)

			return "", &pserror.PasswordSharingError{
				Code:    pserror.DbCommandError,
				Message: message,
			}
		}

		s.log.Debug("link generated")
		return link, nil
	}
}

const recordNotFoundError = "record not found"

func (s *passwordService) GetPasswordFromLink(c context.Context, link string) (string, error) {
	db, err := s.dbFactory.InitDB(c)
	if err != nil {
		return "", initDbError(s.log)
	}

	result := &model.Password{}
	query := db.Where(&model.Password{Link: link}).First(result)
	if err := query.Error; err != nil {
		if err.Error() == recordNotFoundError {
			const message = "password not found"

			s.log.Warn(message,
				zap.String("link", link),
			)

			return "", &pserror.PasswordSharingError{
				Code:    pserror.PasswordNotFound,
				Message: message,
			}
		}

		const message = "error on db query"

		s.log.Error(message,
			zap.Error(err),
		)

		return "", &pserror.PasswordSharingError{
			Code:    pserror.RandomizerError,
			Message: message,
		}
	}

	return result.Password, nil
}

func initDbError(log *zap.Logger) error {
	const message = "failed to init db"

	log.Error(message)

	return &pserror.PasswordSharingError{
		Code:    pserror.InitDbError,
		Message: message,
	}
}
