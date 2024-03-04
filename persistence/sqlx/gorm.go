package sqlx

import (
	"net/url"
	"strings"
	"time"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/sqlx/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGorm(conf *config.Config, logger logging.Logger) (*gorm.DB, error) {
	options := &gorm.Config{
		Logger: NewLogger(logger),
	}

	if u, err := url.ParseRequestURI(conf.Uri); err == nil {
		// GORM perform write (create/update/delete) operations run inside a transaction to ensure data consistency,
		// you can disable it during initialization if it is not required,
		// you will gain about 30%+ performance improvement after that
		options.SkipDefaultTransaction = u.Query().Get("skip_default_transaction") != ""
	}

	var orm *gorm.DB
	var err error

	if strings.HasPrefix(conf.Uri, "postgres") {
		orm, err = gorm.Open(postgres.Open(conf.Uri), options)
	} else {
		orm, err = gorm.Open(sqlite.Open(conf.Uri), options)
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
