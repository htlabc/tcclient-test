package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	"githup.com/htl/tcclienttest/internal/appserver/store"
	"githup.com/htl/tcclienttest/internal/pkg/options"
	"githup.com/htl/tcclienttest/pkg/db"
	"gorm.io/gorm"
	"sync"
)

type datastore struct {
	db *gorm.DB

	// can include two database instance if needed
	// docker *grom.DB
	// db *gorm.DB
}

func (ds *datastore) Images() store.ImageStore {
	return newImages(ds)
}

func (ds *datastore) Terminals() store.TerminalStore {
	return newTerminals(ds)
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed ")
	}
	return db.Close()
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

//创建mysql工厂
// GetMySQLFactoryOr create mysql factory with the given config.
func GetMySQLFactoryOr(opts *options.MySQLOptions) (store.Factory, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              opts.LogLevel,
		}
		dbIns, err = db.New(options)

		// uncomment the following line if you need auto migration the given models
		// not suggested in production environment.
		// migrateDatabase(dbIns)

		mysqlFactory = &datastore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	return mysqlFactory, nil
}
