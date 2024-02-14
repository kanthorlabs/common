package sqlx

import (
	"strings"
	"time"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/sqlx/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGorm(conf *config.Config, logger logging.Logger) (*gorm.DB, error) {
	opts := &gorm.Config{
		// GORM perform write (create/update/delete) operations run inside a transaction to ensure data consistency,
		// you can disable it during initialization if it is not required,
		// you will gain about 30%+ performance improvement after that
		SkipDefaultTransaction: conf.SkipDefaultTxn,
		Logger:                 NewLogger(logger),
	}

	var orm *gorm.DB
	var err error

	if strings.HasPrefix(conf.Uri, "postgres") {
		orm, err = gorm.Open(postgres.Open(conf.Uri), opts)
	} else {
		orm, err = gorm.Open(sqlite.Open(conf.Uri), opts)
	}
	if err != nil {
		return nil, err
	}

	db, err := orm.DB()
	if err != nil {
		return nil, err
	}

	// each connection has their backend
	// the longer the connection is alive, the more memory they consume
	db.SetConnMaxLifetime(time.Millisecond * time.Duration(conf.Connection.MaxLifetime))
	db.SetConnMaxIdleTime(time.Millisecond * time.Duration(conf.Connection.MaxIdletime))
	db.SetMaxIdleConns(conf.Connection.MaxIdleCount)
	db.SetMaxOpenConns(conf.Connection.MaxOpenCount)

	return orm, nil
}
