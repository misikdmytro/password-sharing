package service

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/database"
	pserror "github.com/misikdmitriy/password-sharing/error"
	"github.com/misikdmitriy/password-sharing/helper"
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

func (s *passwordService) CreateLinkFromPassword(c context.Context, pwd string) (string, error) {
	db, close, err := s.dbFactory.InitDB(c)
	if err != nil {
		return "", initDbError(s.log)
	}
	defer close()

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

		var command *gorm.DB
		measureTime(func() {
			command = db.Save(model.NewPassword(link, pwd))
		}, dbTime.WithLabelValues(newPassword))
		dbCounter.WithLabelValues(newPassword).Inc()

		if err := command.Error; err != nil {
			pgErr, ok := err.(*pgconn.PgError)
			if ok && pgErr.Code == pgUniqueViolationCode {
				dbErrorsCounter.WithLabelValues(uniqueViolation).Inc()
				s.log.Warn("retry after unique key violation")
				continue
			}

			const message = "error on db command"

			dbErrorsCounter.WithLabelValues(unknownError).Inc()
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
	db, close, err := s.dbFactory.InitDB(c)
	if err != nil {
		return "", initDbError(s.log)
	}
	defer close()

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
			s.log.Warn(message,
				zap.String("link", link),
			)

			return "", &pserror.PasswordSharingError{
				Code:    pserror.PasswordNotFound,
				Message: message,
			}
		}

		const message = "error on db query"

		dbErrorsCounter.WithLabelValues(unknownError).Inc()
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
