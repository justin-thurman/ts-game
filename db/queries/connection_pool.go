package queries

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context) (*postgres, error) {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL env var not set")
	}
	var outerErr error
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, dbUrl)
		if err != nil {
			outerErr = fmt.Errorf("unable to create connection pool: %v", err)
			return
		}
		pgInstance = &postgres{db}
	})
	if outerErr != nil {
		return nil, outerErr
	}
	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func (pg *postgres) Exec(ctx context.Context, sqlStr string, args ...interface{}) (pgconn.CommandTag, error) {
	return pg.db.Exec(ctx, sqlStr, args...)
}

func (pg *postgres) Query(ctx context.Context, sqlString string, args ...interface{}) (pgx.Rows, error) {
	return pg.db.Query(ctx, sqlString, args...)
}

func (pg *postgres) QueryRow(ctx context.Context, sqlString string, args ...interface{}) pgx.Row {
	return pg.db.QueryRow(ctx, sqlString, args...)
}
