package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type postgresContainer struct {
	container testcontainers.Container
	URI       string
	Terminate func(context.Context) error
}

func startPostgresContainer(ctx context.Context) (*postgresContainer, error) {
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

	return &postgresContainer{container, connectionString, container.Terminate}, nil
}

func TestSchemaUp(t *testing.T) {
	ctx := context.Background()
	container, err := startPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	log.Println("Connection URI:", container.URI)

	os.Args = []string{"", "up"}
	os.Setenv("DATABASE_URL", container.URI)
	main()
}

func TestSchemaDown(t *testing.T) {
	ctx := context.Background()
	container, err := startPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	log.Println("Connection URI:", container.URI)

	os.Args = []string{"", "down"}
	os.Setenv("DATABASE_URL", container.URI)
	main()
}

func TestSchemaReset(t *testing.T) {
	ctx := context.Background()
	container, err := startPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	log.Println("Connection URI:", container.URI)

	os.Args = []string{"", "reset"}
	os.Setenv("DATABASE_URL", container.URI)
	main()
}
