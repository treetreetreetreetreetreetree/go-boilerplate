package database

import (
	"database/sql"
	"errors"
	"fmt"
	"go-boilerplate/config"
	"strconv"

	"log/slog"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ErrDBConnectionFailed = errors.New("db connection was failed")

type Database struct {
	Gorm *gorm.DB
	SQL  *sql.DB
}

func Setup(cfg *config.DatabaseConfig) (*Database, error) {
	var err error

	p, err := strconv.Atoi(cfg.Port)
	if err != nil {
		slog.Error("[SQL] parse config", "error", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", cfg.Host, cfg.User, cfg.Password, cfg.Name, p)

	var gormCfg = &gorm.Config{}

	if cfg.Debug > 0 {
		gormCfg = &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)

	if db.Error != nil || err != nil {
		slog.Error("[SQL] ErrorDBConnectionFailed", "error", err)
		return nil, ErrDBConnectionFailed
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("[SQL] ErrorDBConnectionFailed", "error", err)
		return nil, ErrDBConnectionFailed
	}

	var ping bool
	db.Raw("select 1").Scan(&ping)
	if !ping {
		slog.Error("[SQL] db connection was failed", "error", err)
		return nil, ErrDBConnectionFailed
	}

	slog.Info("[SQL]", "message", "connection was successfully opened to database")

	return &Database{Gorm: db, SQL: sqlDB}, nil
}

func (d *Database) EnsureMigrations(migrations []*gormigrate.Migration) error {
	m := gormigrate.New(d.Gorm, &gormigrate.Options{
		TableName:                 "gorm_migrations",
		IDColumnName:              "id",
		IDColumnSize:              512,
		UseTransaction:            false,
		ValidateUnknownMigrations: false,
	}, migrations)

	if err := m.Migrate(); err != nil {
		slog.Error("[SQL] could not migrate", "err", err)
		return err
	}

	slog.Info("[SQL]", "message", "migration did run successfully")
	return nil
}
