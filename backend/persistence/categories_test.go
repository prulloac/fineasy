package persistence

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var p = Persistence{}

func TestInsertCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	category := Category{Name: "Food", Icon: "restaurant", Color: "red", Description: "Food and drinks"}
	mock.ExpectQuery("SELECT id FROM categories").WithArgs(category.Name).WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO categories").WithArgs(category.Name, category.Icon, category.Color, category.Description).WillReturnResult(sqlmock.NewResult(1, 1))

	p.db = db
	err = p.InsertCategory(category)

	if err != nil {
		t.Errorf("error was not expected while inserting category: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	category := Category{Name: "Food", Icon: "restaurant", Color: "red", Description: "Food and drinks"}
	mock.ExpectQuery("SELECT id, name, icon, color, description FROM categories").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "icon", "color", "description"}).AddRow(1, category.Name, category.Icon, category.Color, category.Description))

	p.db = db
	r, err := p.GetCategories()

	if err != nil {
		t.Errorf("error was not expected while getting categories: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, c := range r {
		if c.Name != category.Name {
			t.Errorf("expected %s but got %s", category.Name, c.Name)
		}
	}
}

func TestGetCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	category := Category{ID: 1, Name: "Food", Icon: "restaurant", Color: "red", Description: "Food and drinks"}
	mock.ExpectQuery("SELECT id, name, icon, color, description, deleted_at, created_at, updated_at FROM categories").
		WithArgs(category.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "icon", "color", "description", "deleted_at", "created_at", "updated_at"}).
			AddRow(1, category.Name, category.Icon, category.Color, category.Description, category.DeletedAt, category.CreatedAt, category.UpdatedAt),
		)

	p.db = db
	r, err := p.GetCategory(1)

	if err != nil {
		t.Errorf("error was not expected while getting category: %s", err)
	}

	if r.Name != category.Name {
		t.Errorf("expected %s but got %s", category.Name, r.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	category := Category{ID: 1, Name: "Food", Icon: "restaurant", Color: "red", Description: "Food and drinks"}
	mock.ExpectExec("UPDATE categories").WithArgs(category.Name, category.Icon, category.Color, category.Description, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	p.db = db
	err = p.UpdateCategory(category)

	if err != nil {
		t.Errorf("error was not expected while updating category: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectExec("UPDATE categories").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	p.db = db
	err = p.DeleteCategory(1)

	if err != nil {
		t.Errorf("error was not expected while deleting category: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRestoreCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	category := Category{ID: 1, Name: "Food", Icon: "restaurant", Color: "red", Description: "Food and drinks"}

	mock.ExpectExec("UPDATE categories").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT id, name, icon, color, description, deleted_at, created_at, updated_at FROM categories").
		WithArgs(category.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "icon", "color", "description", "deleted_at", "created_at", "updated_at"}).
			AddRow(1, category.Name, category.Icon, category.Color, category.Description, category.DeletedAt, category.CreatedAt, category.UpdatedAt),
		)

	p.db = db
	c, err := p.RestoreCategory(1)

	if err != nil {
		t.Errorf("error was not expected while restoring category: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if c.ID != 1 {
		t.Errorf("expected 1 but got %d", c.ID)
	}
}
