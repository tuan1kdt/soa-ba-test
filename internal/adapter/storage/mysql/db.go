package mysql

import (
	"context"
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/tuan1kdt/soa-ba-test/internal/adapter/config"
	"gorm.io/gorm"

	"github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
)

type DB struct {
	*gorm.DB
	url string
}

// New creates a new PostgreSQL database instance
func New(ctx context.Context, config *config.DB) (*DB, error) {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		config.Connection,
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	connConfig := gormmysql.Config{
		DriverName: "mysql",
		DSN:        url,
		DSNConfig:  &mysql.Config{},
	}

	dialector := gormmysql.New(connConfig)

	db, err := gorm.Open(dialector)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxIdleConns)
	}
	if config.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(config.MaxLifetime)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("cannot ping database, err: '%v'", err)
	}

	return &DB{
		db,
		url,
	}, nil
}

// Migrate runs the database migration
func (db *DB) Migrate() error {
	return nil
}

// ErrorCode returns the error code of the given error
func (db *DB) ErrorCode(err error) string {
	return err.Error()
}

// Close closes the database connection
func (db *DB) Close() {
}
