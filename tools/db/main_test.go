package main

import (
	"context"
	"os"
	"testing"

	"github.com/prulloac/fineasy/tests"
)

func TestSchemaUp(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)

	os.Args = []string{"", "up"}
	os.Setenv("DATABASE_URL", container.URI)
	main()
	container.Terminate(ctx)
}

func TestSchemaDown(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)

	os.Args = []string{"", "down"}
	os.Setenv("DATABASE_URL", container.URI)
	main()
	container.Terminate(ctx)
}

func TestSchemaReset(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)

	os.Args = []string{"", "reset"}
	os.Setenv("DATABASE_URL", container.URI)
	main()
	container.Terminate(ctx)
}
