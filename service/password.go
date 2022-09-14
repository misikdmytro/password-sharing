package service

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/database"
	pserror "github.com/misikdmitriy/password-sharing/error"
	"github.com/misikdmitriy/password-sharing/helper"
	"github.com/misikdmitriy/password-sharing/logger"
	"github.com/misikdmitriy/password-sharing/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PasswordService interface {
	GetPasswordFromLink(context.Context, string) (string, error)
	CreateLinkFromPassword(context.Context, string) (string, error)
}

type passwordService struct {
	dbFactory     database.DbFactory
	configuration *config.Config
	randomFactory helper.RandomGeneratorFactory
	loggerFactory logger.LoggerFactory
	encoder       helper.Encoder
}

func NewPasswordService(dbFactory database.DbFactory,
	conf *config.Config,
	rf helper.RandomGeneratorFactory,
	loggerFactory logger.LoggerFactory,
	encoder helper.Encoder) PasswordService {
	return &passwordService{
		dbFactory:     dbFactory,
		configuration: conf,
		randomFactory: rf,
		loggerFactory: loggerFactory,
		encoder:       encoder,
	}
}

const pgUniqueViolationCode = "23505"

var (
	dbCounter *prometheus.CounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "password_sharing_db",
		Help: "The total number of DB queries/commands",
	}, []string{"type"})

	dbErrorsCounter *prometheus.CounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "password_sharing_db_errors",
		Help: "The total number of errors on DB queries",
	}, []string{"type"})

	dbTime *prometheus.HistogramVec = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "password_sharing_db_time",
		Help:    "The time of DB queries/commands",
		Buckets: prometheus.DefBuckets,
	}, []string{"type"})
)

const (
	newPassword     = "new_password"
	getPassword     = "get_password"
	uniqueViolation = "unique_violation"
	unknownError    = "unknown_error"
	notFound        = "not_found"
)

func (s *passwordService) CreateLinkFromPassword(c context.Context, password string) (string, error) {
	appLogger, loggerClose, err := s.loggerFactory.NewLogger()
	if err != nil {
		return "", err
	}
	defer loggerClose()

	db, dbClose, err := s.dbFactory.InitDB(c)
	if err != nil {
		return "", initDbError(appLogger)
	}
	defer dbClose()

	encoded, err := s.encoder.Encode(password)
	if err != nil {
		const message = "failed on encoding"

		appLogger.Error(message)

		return "", &pserror.PasswordSharingError{
			Code:    pserror.EncodeError,
			Message: message,
		}
	}

	for {
		rg := s.randomFactory.NewRandomGenerator()
		link, err := rg.RandomString(s.configuration.App.LinkLength)
		if err != nil {
			const message = "error on randomizing"

			appLogger.Error(message,
				zap.Error(err),
				zap.Int("length", s.configuration.App.LinkLength),
			)

			return "", &pserror.PasswordSharingError{
				Code:    pserror.RandomizerError,
				Message: message,
			}
		}

		var command *gorm.DB
		measureTime(func() {
			command = db.Save(&model.Password{
				Link:     link,
				Password: encoded,
			})
		}, dbTime.WithLabelValues(newPassword))
		dbCounter.WithLabelValues(newPassword).Inc()

		if err := command.Error; err != nil {
			pgErr, ok := err.(*pgconn.PgError)
			if ok && pgErr.Code == pgUniqueViolationCode {
				dbErrorsCounter.WithLabelValues(uniqueViolation).Inc()
				appLogger.Warn("retry after unique key violation")
				continue
			}

			const message = "error on db command"

			dbErrorsCounter.WithLabelValues(unknownError).Inc()
			appLogger.Error(message,
				zap.Error(err),
			)

			return "", &pserror.PasswordSharingError{
				Code:    pserror.DbCommandError,
				Message: message,
			}
		}

		appLogger.Debug("link generated")
		return link, nil
	}
}

const recordNotFoundError = "record not found"

func (s *passwordService) GetPasswordFromLink(c context.Context, link string) (string, error) {
	appLogger, loggerClose, err := s.loggerFactory.NewLogger()
	if err != nil {
		return "", err
	}
	defer loggerClose()

	db, dbClose, err := s.dbFactory.InitDB(c)
	if err != nil {
		return "", initDbError(appLogger)
	}
	defer dbClose()

	result := &model.Password{}
	var query *gorm.DB
	measureTime(func() {
		query = db.Where(&model.Password{Link: link}).First(result)
	}, dbTime.WithLabelValues(getPassword))
	dbCounter.WithLabelValues(getPassword).Inc()

	if err := query.Error; err != nil {
		if err.Error() == recordNotFoundError {
			const message = "password not found"

			dbErrorsCounter.WithLabelValues(notFound).Inc()
			appLogger.Warn(message,
				zap.String("link", link),
			)

			return "", &pserror.PasswordSharingError{
				Code:    pserror.PasswordNotFound,
				Message: message,
			}
		}

		const message = "error on db query"

		dbErrorsCounter.WithLabelValues(unknownError).Inc()
		appLogger.Error(message,
			zap.Error(err),
		)

		return "", &pserror.PasswordSharingError{
			Code:    pserror.RandomizerError,
			Message: message,
		}
	}

	decoded, err := s.encoder.Decode(result.Password)
	if err != nil {
		const message = "failed on decoding"

		appLogger.Error(message)

		return "", &pserror.PasswordSharingError{
			Code:    pserror.DecodeError,
			Message: message,
		}
	}

	return decoded, nil
}

func measureTime(action func(), metric prometheus.Observer) {
	timer := prometheus.NewTimer(metric)
	action()
	timer.ObserveDuration()
}

func initDbError(log *zap.Logger) error {
	const message = "failed to init db"

	log.Error(message)

	return &pserror.PasswordSharingError{
		Code:    pserror.InitDbError,
		Message: message,
	}
}
