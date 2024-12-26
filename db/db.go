package db

import (
	"backendsetup/m/config"
	"backendsetup/m/db/sql/dbgen"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Init(conf *config.Config) *dbgen.Queries {
	ctx := context.Background()
	pgString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort, conf.DBName)
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

	m, err := migrate.New("file://db/sql/migrations/", pgString)
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
