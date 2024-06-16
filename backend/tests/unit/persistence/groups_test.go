package persistence__test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/prulloac/fineasy/persistence"
	. "github.com/prulloac/fineasy/persistence/entity"
)

func TestInsertGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sampleTime, err := time.Parse(time.DateOnly, "2021-01-01")

	group := Group{Name: "Family", CreatedBy: 1, CreatedAt: sampleTime}
	mock.ExpectQuery("SELECT id FROM groups").
		WithArgs(group.Name).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO groups").
		WithArgs(group.Name, group.CreatedBy).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = persistence.NewPersistence(db)
	err = p.GetGroupRepository().InsertGroup(group)

	if err != nil {
		t.Errorf("error was not expected while inserting group: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetGroups(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sampleTime, err := time.Parse(time.DateOnly, "2021-01-01")

	group := Group{Name: "Family", CreatedBy: 1, CreatedAt: sampleTime}
	mock.ExpectQuery("SELECT id, name, created_by, created_at FROM groups").
		WithArgs(group.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_by", "created_at"}).
			AddRow(1, group.Name, group.CreatedBy, group.CreatedAt))

	var p = persistence.NewPersistence(db)
	r, err := p.GetGroupRepository().GetGroups(group.CreatedBy)

	if err != nil {
		t.Errorf("error was not expected while getting groups: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, g := range r {
		if g.Name != group.Name {
			t.Errorf("expected: %s, got: %s", group.Name, g.Name)
		}
	}
}

func TestGetGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sampleTime, err := time.Parse(time.DateOnly, "2021-01-01")

	group := Group{Name: "Family", CreatedBy: 1, CreatedAt: sampleTime}
	mock.ExpectQuery("SELECT id, name, created_by, created_at FROM groups").
		WithArgs(group.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_by", "created_at"}).
			AddRow(1, group.Name, group.CreatedBy, group.CreatedAt))

	var p = persistence.NewPersistence(db)
	r, err := p.GetGroupRepository().GetGroup(group.ID)

	if err != nil {
		t.Errorf("error was not expected while getting group: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if r.Name != group.Name {
		t.Errorf("expected: %s, got: %s", group.Name, r.Name)
	}
}

func TestUpdateGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	group := Group{ID: 1, Name: "Family"}
	mock.ExpectExec("UPDATE groups").
		WithArgs(group.Name, group.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = persistence.NewPersistence(db)
	err = p.GetGroupRepository().UpdateGroup(group)

	if err != nil {
		t.Errorf("error was not expected while updating group: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
