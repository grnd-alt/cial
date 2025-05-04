package db

import (
	"backendsetup/m/db/sql/dbgen"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(username string, password string, host string, port int, dbname string) *dbgen.Queries {
	ctx := context.Background()
	pgString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbname)
	var conn *pgxpool.Pool
	var err error
	for {
		conn, err = pgxpool.New(ctx, pgString)
		if err != nil {
			log.Fatal(err)
		}

		err = conn.Ping(ctx)
		if err != nil {
			log.Println("DB conn failed retrying in 5s")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	migrationPath := fmt.Sprintf("file://%s", filepath.Join(basePath, "db/sql/migrations/"))

	m, err := migrate.New(migrationPath, pgString)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		if err.Error() == "no change" {
			fmt.Println("No changes to migrate")
		} else {
			log.Panic(err)
		}
	}

	queries := dbgen.New(conn)

	return queries
}
