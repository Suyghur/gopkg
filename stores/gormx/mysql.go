//@File     mysql.go
//@Time     2024/7/18
//@Author   #Suyghur,

package gormx

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormDB(dsn string, opts ...Option) *gorm.DB {
	c := &config{
		MaxIdleConns:    defaultMaxIdleConns,
		MaxOpenConns:    defaultMaxOpenConns,
		ConnMaxLifetime: defaultConnMaxLifetime,
		SlowThreshold:   defaultSlowThreshold,
		DialTimeout:     defaultDialTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	gormDB, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db, _ := gormDB.DB()

	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)

	registerHook(gormDB, durationHook(c.SlowThreshold), traceHook())

	return gormDB
}

func NewGormDBFromConn(conn *sql.DB, opts ...Option) *gorm.DB {
	c := &config{
		MaxIdleConns:    defaultMaxIdleConns,
		MaxOpenConns:    defaultMaxOpenConns,
		ConnMaxLifetime: defaultConnMaxLifetime,
		SlowThreshold:   defaultSlowThreshold,
		DialTimeout:     defaultDialTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      conn,
	}), &gorm.Config{
		Logger: NewLogger(),
	})

	db, _ := gormDB.DB()

	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)

	registerHook(gormDB, durationHook(c.SlowThreshold), traceHook())

	return gormDB
}
