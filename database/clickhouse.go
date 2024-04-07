package database

import (
	"fmt"
	"log"

	"github.com/cetnfurkan/core/config"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type (
	clickhouseDatabase struct {
		cfg    *clickhouseDatabaseConfig
		client *gorm.DB
	}

	clickhouseDatabaseConfig struct {
		*config.Database
	}
)

// NewClickhouseDatabase creates a new clickhouse database instance.
//
// It takes a config instance and returns a new database interface instance.
//
// It will panic
// if it fails to connect to clickhouse database.
func NewClickhouseDatabase(cfg *config.Database, opts ...option[gorm.DB]) Database {
	var (
		err error
	)

	database := &clickhouseDatabase{
		cfg: &clickhouseDatabaseConfig{
			Database: cfg,
		},
	}

	database.client, err = gorm.Open(clickhouse.Open(database.getDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to clickhouse database: %v", err)
	}

	for _, opt := range opts {
		err = opt(database.client)
		if err != nil {
			log.Fatalf("failed to apply option: %v", err)
		}
	}

	return database
}

func (database *clickhouseDatabase) UnmarshalExtra() {}

func (database *clickhouseDatabase) Get() any {
	return database.client
}

func (database *clickhouseDatabase) getDSN() string {
	return fmt.Sprintf(
		"clickhouse://%s:%s@%s:%d/%s?dial_timeout=10s&read_timeout=20s",
		database.cfg.Database.User,
		database.cfg.Database.Password,
		database.cfg.Database.Host,
		database.cfg.Database.Port,
		database.cfg.Database.Name,
	)
}
