package persistence

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Test() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL env var not set")
	}
	dbpool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "SELECT 'Hello, world!'").Scan(&greeting)
	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}
	log.Println(greeting)

	query := `INSERT INTO foo (id, name) VALUES (@id, @name)`
	args := pgx.NamedArgs{
		"id":   1,
		"name": "Justin",
	}
	ct, err := dbpool.Exec(context.Background(), query, args)
	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}
	log.Println(ct.String())

	rows, err := dbpool.Query(context.Background(), "SELECT * FROM foo")
	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatalf("Error reading row values: %v\n", err)
		}
		log.Println(values)
	}

	log.Println(greeting)
}
