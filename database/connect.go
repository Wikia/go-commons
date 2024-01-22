package database

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"gorm.io/plugin/opentelemetry/tracing"
	"moul.io/zapgorm2"
)

func GetConnection(logger *zap.Logger, logLevel dblogger.LogLevel, sources, replicas []string, connMaxIdleTime, connMaxLifeTime time.Duration, maxIdleConns, maxOpenConns int) (*gorm.DB, error) {
	dbLogger := zapgorm2.New(logger)
	dbLogger.SetAsDefault()

	if len(sources) == 0 {
		return nil, errors.New("no source database provided")
	}

	db, err := gorm.Open(mysql.Open(sources[0]), &gorm.Config{Logger: dbLogger.LogMode(logLevel)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	dbSources := make([]gorm.Dialector, len(sources))
	for _, dsn := range sources {
		dbSources = append(dbSources, mysql.Open(dsn))
	}

	dbReplicas := make([]gorm.Dialector, len(replicas))
	for _, dsn := range replicas {
		dbReplicas = append(dbReplicas, mysql.Open(dsn))
	}

	err = db.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:  dbSources,
			Replicas: dbReplicas,
			Policy:   dbresolver.RandomPolicy{},
		}).
			SetConnMaxIdleTime(connMaxIdleTime).
			SetConnMaxLifetime(connMaxLifeTime).
			SetMaxIdleConns(maxIdleConns).
			SetMaxOpenConns(maxOpenConns))

	if err != nil {
		return nil, errors.Wrap(err, "could not setup db replicas")
	}

	if err = db.Use(tracing.NewPlugin()); err != nil {
		return nil, errors.Wrap(err, "could not setup db tracing and metrics")
	}

	return db, nil
}
