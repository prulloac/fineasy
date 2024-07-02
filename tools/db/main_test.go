package main

import (
	"context"
	"os"
	"testing"

	"github.com/prulloac/fineasy/tests"
)

func TestSchemaUp(t *testing.T) {
	ctx := context.Background()
	container, err := tests.StartPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	os.Args = []string{"", "up"}
	os.Setenv("DATABASE_URL", container.URI)
	main()
}

func TestSchemaDown(t *testing.T) {
	ctx := context.Background()
	container, err := tests.StartPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	os.Args = []string{"", "down"}
	os.Setenv("DATABASE_URL", container.URI)
	main()
}

func TestSchemaReset(t *testing.T) {
	ctx := context.Background()
	container, err := tests.StartPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	os.Args = []string{"", "reset"}
	os.Setenv("DATABASE_URL", container.URI)
	main()
}
