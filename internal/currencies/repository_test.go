package currencies

import (
	"context"
	"testing"

	"github.com/prulloac/fineasy/tests"
)

func TestCreateAndDropTables(t *testing.T) {
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)

	var p = CurrencyRepository{DB: container.DB}

	err := p.CreateTable()
	if err != nil {
		t.Errorf("error was not expected while creating tables: %s", err)
	}

	err = p.DropTable()
	if err != nil {
		t.Errorf("error was not expected while dropping tables: %s", err)
	}

	container.Terminate(ctx)
}
