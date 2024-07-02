package tests

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	container testcontainers.Container
	URI       string
	Terminate func(context.Context) error
	DB        func() (*sql.DB, error)
}

func StartPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:alpine"),
		postgres.WithDatabase("fineasy"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		return nil, err
	}

	state, err := container.State(ctx)

	if err != nil {
		return nil, err
	}

	fmt.Println("Container state:", state.Running)

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		return nil, err
	}

	return &PostgresContainer{container, connectionString, container.Terminate, func() (*sql.DB, error) { return sql.Open("postgres", connectionString) }}, nil
}
