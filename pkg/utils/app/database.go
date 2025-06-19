package app

import (
	"context"
	"fmt"
	"log"
	"time"

	dbCon "example.com/m/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db   *dbCon.Queries
	Db   *dbCon.Queries
	Conn *pgx.Conn
)

// Connect connect to the database
func Connect(config AppConfig) (*dbCon.Queries, *dbCon.Store, *pgxpool.Pool) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DATABASE_USERNAME,
		config.DATABASE_PASSWORD,
		config.DATABASE_HOST,
		config.DATABASE_PORT,
		config.DATABASE_NAME)

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
	}

	poolConfig.MaxConns = 20
	poolConfig.MinConns = 1
	poolConfig.HealthCheckPeriod = 5 * time.Second

	// Create a new connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	store := dbCon.NewStore(pool)
	db = dbCon.New(pool)

	return db, store, pool
}
