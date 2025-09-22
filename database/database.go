package database

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

type DBType int

const (
	Postgres DBType = iota
	Mysql
	Sqlite
)

type database struct {
	path   string
	host   string
	port   int
	name   string
	user   string
	passwd string
}

func NewDatabase(dbType DBType, log *logrus.Logger, opts ...Option) *gorm.DB {
	db := &database{}

	// set default
	switch dbType {
	case Postgres:
		db.port = 5432
	case Mysql:
		db.port = 3306
	}

	// apply options
	for _, opt := range opts {
		opt(db)
	}

	dialector, err := db.buildGormDialector(dbType)
	if err != nil {
		log.Fatalln(err)
	}

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlog.New(log, gormlog.Config{
			SlowThreshold: 200 * time.Millisecond, // Slow SQL threshold
			Colorful:      true,                   // enable color
		}),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalln("Can't connect to database")
	}

	return gormDB
}

func (db *database) buildGormDialector(dbType DBType) (gorm.Dialector, error) {
	switch dbType {
	case Mysql:
		dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s",
			db.host,
			db.user,
			db.passwd,
			db.port,
			db.name,
		)
		return mysql.Open(dsn), nil
	case Postgres:
		dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s",
			db.host,
			db.user,
			db.passwd,
			db.port,
			db.name,
		)
		return postgres.Open(dsn), nil
	case Sqlite:
		dsn := fmt.Sprintf("%s",
			db.path,
		)
		return sqlite.Open(dsn), nil
	default:
		return nil, fmt.Errorf("Unsupport DB type: %v", dbType)
	}
}

// NOTE: Options

type Option func(*database)

func WithPath(path string) Option {
	return func(db *database) {
		db.path = path
	}
}

func WithHost(host string) Option {
	return func(db *database) {
		db.host = host
	}
}

func WithPort(port int) Option {
	return func(db *database) {
		db.port = port
	}
}

func WithName(name string) Option {
	return func(db *database) {
		db.name = name
	}
}

func WithUser(user string) Option {
	return func(db *database) {
		db.user = user
	}
}

func WithPasswd(passwd string) Option {
	return func(db *database) {
		db.passwd = passwd
	}
}
