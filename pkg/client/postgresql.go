package client

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
	"weather_service/internal/config"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

// NewClient creates a new DB client for PostgreSQL server using pgx package
func NewClient(ctx context.Context, cfg config.Config) (pool *pgxpool.Pool, err error) {
	// connect string
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	//Create connection pool
	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("Can not connect DB: error", err)
		return nil, err
	}

	log.Println("DB client connected successfully")

	// migrate DB
	migrate, err := os.ReadFile("migrate.sql")
	if err != nil {
		log.Fatal("Can not read migrate.sql: error", err)
		return nil, err
	}
	_, err = pool.Exec(ctx, string(migrate))
	if err != nil {
		log.Fatal("Can not migrate DB: error", err)
		return nil, err
	}
	log.Println("DB migrated successfully")

	// ping
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal("Can not ping DB: error", err)
		return nil, err
	}
	log.Println("DB pinged successfully")

	return pool, nil
}
