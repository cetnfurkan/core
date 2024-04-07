package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/cetnfurkan/core/config"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type (
	postgresDatabase[T any] struct {
		cfg    *postgresDatabaseConfig
		client *T
		pool   *sql.DB
	}

	postgresDatabaseConfig struct {
		*config.Database
		Extra postgresDatabaseConfigExtra
	}

	postgresDatabaseConfigExtra struct {
		SSLMode                  string `mapstructure:"sslmode"`
		ConnRetry                int    `mapstructure:"connRetry"`
		MaxOpenConns             int    `mapstructure:"maxOpenConns"`
		MaxIdleConns             int    `mapstructure:"maxIdleConns"`
		DescriptionCacheCapacity int    `mapstructure:"descriptionCacheCapacity"`
		StatementCacheCapacity   int    `mapstructure:"statementCacheCapacity"`
		ConnTimeOut              int    `mapstructure:"connTimeOut"`
		MaxOpenConnTTL           int    `mapstructure:"maxOpenConnTTL"`
		MaxIdleConnTTL           int    `mapstructure:"maxIdleConnTTL"`
		QueryExecMode            string `mapstructure:"queryExecMode"`
	}

	ConnPool interface {
		Acquire(ctx context.Context) (*pgxpool.Conn, error)
	}
)

// NewPostgresDatabase creates a new postgres database instance.
//
// It takes a config instance and returns a new database interface instance.
//
// It will panic
// if it fails to unmarhal extra config data,
// if it fails to create a new PGX pool or
// if it fails to connect to postgres database.
func NewPostgresDatabase[T any](cfg *config.Database, createClient func(*entsql.Driver) *T, opts ...option[T]) Database {
	database := &postgresDatabase[T]{
		cfg: &postgresDatabaseConfig{
			Database: cfg,
		},
	}

	database.UnmarshalExtra()

	err := database.createPool()
	if err != nil {
		log.Fatal("failed to create new PGX: ", err)
	}

	driver := entsql.OpenDB(dialect.Postgres, database.pool)
	database.client = createClient(driver)

	for _, opt := range opts {
		err := opt(database.client)
		if err != nil {
			log.Fatalf("failed to apply option: %v", err)
		}
	}

	return database
}

func (database *postgresDatabase[T]) Get() any {
	return database.client
}

func (database *postgresDatabase[T]) UnmarshalExtra() {
	err := mapstructure.Decode(database.cfg.Database.Extra, &database.cfg.Extra)
	if err != nil {
		log.Fatal("failed to decode extra config into struct: ", err)
	}
}

func (database *postgresDatabase[T]) createPool() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(database.cfg.Extra.ConnTimeOut)*time.Second)
	defer cancel()

	database.pool = stdlib.OpenDB(database.getConnectionConfig())

	database.pool.SetConnMaxIdleTime(time.Duration(database.cfg.Extra.MaxIdleConnTTL) * time.Second)
	database.pool.SetConnMaxLifetime(time.Duration(database.cfg.Extra.MaxOpenConnTTL) * time.Second)
	database.pool.SetMaxOpenConns(database.cfg.Extra.MaxOpenConns)
	database.pool.SetMaxIdleConns(database.cfg.Extra.MaxIdleConns)

	err := database.pingConnection(ctx, database.pool)
	if err != nil {
		return errors.Wrap(err, "postgres connection pool ping failed")
	}

	return nil
}

func (database *postgresDatabase[T]) getConnectionConfig() pgx.ConnConfig {
	connectionConfig, err := pgx.ParseConfig(database.getDSN())
	if err != nil {
		log.Fatal("failed to parse connection config: ", err)
	}

	return *connectionConfig
}

func (database *postgresDatabase[T]) getDSN() string {
	var (
		query = make(url.Values)
	)

	if database.cfg.Extra.StatementCacheCapacity >= 0 {
		query.Set("statement_cache_capacity", strconv.Itoa(database.cfg.Extra.StatementCacheCapacity))
	}

	if database.cfg.Extra.DescriptionCacheCapacity >= 0 {
		query.Set("description_cache_capacity", strconv.Itoa(database.cfg.Extra.DescriptionCacheCapacity))
	}

	if database.cfg.Extra.QueryExecMode != "" {
		query.Set("default_query_exec_mode", database.cfg.Extra.QueryExecMode)
	}

	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(database.cfg.User, database.cfg.Password),
		Host:     net.JoinHostPort(database.cfg.Host, fmt.Sprintf("%d", database.cfg.Port)),
		Path:     database.cfg.Name,
		RawQuery: query.Encode(),
	}

	return dsn.String()
}

func (database *postgresDatabase[T]) pingConnection(ctx context.Context, pool *sql.DB) error {
	var (
		err error
	)

	for i := 0; i < database.cfg.Extra.ConnRetry; i++ {

		switch err = pool.PingContext(ctx); err {
		case nil:
			return nil

		case context.Canceled, context.DeadlineExceeded:
			return errors.Wrap(err, "ping database connection timeout exceeded")

		default:
		}
	}

	return errors.Wrap(err, "ping database connection failed")
}
