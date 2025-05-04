package tests

import (
	"backendsetup/m/db"
	"backendsetup/m/db/sql/dbgen"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var Queries *dbgen.Queries
var PostgresContainer testcontainers.Container
var once sync.Once

func SetUp() {
	once.Do(func() {
		fmt.Println("setup running")
		var err error
		PostgresContainer, err = postgres.Run(context.Background(),
			"postgres:17-alpine",
			postgres.WithDatabase("goapp"),
			postgres.WithUsername("postgres"),
			postgres.WithPassword("postgres"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Second)),
		)

		if err != nil {
			log.Printf("failed to initialize postgres: %s", err)
			panic(err)
		}

		port, err := PostgresContainer.MappedPort(context.Background(), "5432")

		if err != nil {
			panic(err)
		}

		Queries = db.Init("postgres", "postgres", "localhost", port.Int(), "goapp")
		err = Queries.InsertUser(context.Background(), dbgen.InsertUserParams{
			Username: "testuser",
			UserID:   "1",
		})
		if err != nil {
			panic(err)
		}
		err = Queries.InsertUser(context.Background(), dbgen.InsertUserParams{
			Username: "testuser2",
			UserID:   "2",
		})
		if err != nil {
			panic(err)
		}
	})
}
