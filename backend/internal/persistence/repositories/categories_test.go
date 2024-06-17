package repositories

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/prulloac/fineasy/internal/persistence/entity"
)

func TestInsertCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	category := entity.Category{Name: "Food", Icon: "restaurant", Color: "red", Description: "Food and drinks", Order: 1, GroupID: 1}
	mock.ExpectQuery("SELECT id FROM categories").
		WithArgs(category.Name).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO categories").
		WithArgs(category.Name, category.Icon, category.Color, category.Description, category.Order, category.GroupID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = CategoriesRepository{db}
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

	category := entity.Category{Name: "Food", Icon: "restaurant", Color: "red", Description: "Food and drinks", Order: 1, GroupID: 1}
	mock.ExpectQuery("SELECT id, name, icon, color, description, ord, group_id FROM categories").
		WithArgs(category.GroupID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "icon", "color", "description", "order", "group_id"}).
			AddRow(1, category.Name, category.Icon, category.Color, category.Description, category.Order, category.GroupID))

	var p = CategoriesRepository{db}
	r, err := p.GetCategories(category.GroupID)

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

	category := entity.Category{ID: 1, Name: "Food", Icon: "restaurant", Color: "red", Description: "Food and drinks", Order: 1, GroupID: 1}
	mock.ExpectQuery("SELECT id, name, icon, color, description, ord, group_id FROM categories").
		WithArgs(category.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "icon", "color", "description", "order", "group_id"}).
			AddRow(1, category.Name, category.Icon, category.Color, category.Description, category.Order, category.GroupID),
		)

	var p = CategoriesRepository{db}
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

	category := entity.Category{ID: 1, Name: "Food", Icon: "restaurant", Color: "red", Description: "Food and drinks", Order: 1}
	mock.ExpectExec("UPDATE categories").
		WithArgs(category.Name, category.Icon, category.Color, category.Description, category.Order, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = CategoriesRepository{db}
	err = p.UpdateCategory(category)

	if err != nil {
		t.Errorf("error was not expected while updating category: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
