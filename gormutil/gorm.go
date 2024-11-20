package gormutil

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	gormzerolog "github.com/gorpher/gone/gormutil/gorm-zerolog"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var defaultDB *gorm.DB

func Init(conf Config, debug bool) {
	db, err := New(debug, conf)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if debug {
		db = db.Debug()
	}
	defaultDB = db
}

func AutoMigrate(models []interface{}) {
	defaultDB.AutoMigrate(models...) // nolint
}

func DB() *gorm.DB {
	return defaultDB
}

type Config struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

func New(debug bool, conf Config) (*gorm.DB, error) {
	var director func(dsn string) gorm.Dialector
	switch conf.Driver {
	case "mysql":
		director = mysql.Open
	case "postgres":
		director = postgres.Open
	case "sqlserver":
		director = sqlserver.Open
	default:
		director = sqlite.Open
	}
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	config.Logger = gormzerolog.NewLoggerInterface(debug, int(200*time.Millisecond), 2)
	db, err := gorm.Open(director(conf.DSN), config)
	if err != nil {
		return nil, err
	}
	if debug {
		db = db.WithContext(log.Logger.WithContext(context.Background()))
	}
	sqlDB, _ := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, err
}
