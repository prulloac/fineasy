package persistence__test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/prulloac/fineasy/persistence"
	. "github.com/prulloac/fineasy/persistence/entity"
)

var sampleTime, err = time.Parse(time.DateOnly, "2021-01-01")

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	user := User{Username: "user", Email: "mail", CreatedAt: sampleTime, UpdateAt: sampleTime}
	mock.ExpectQuery("SELECT id FROM users").
		WithArgs(user.Email).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Username, user.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = persistence.NewPersistence(db)
	err = p.GetUserRepository().InsertUser(user)

	if err != nil {
		t.Errorf("error was not expected while inserting user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sampleTime, err := time.Parse(time.DateOnly, "2021-01-01")

	user := User{Username: "user", Email: "mail", CreatedAt: sampleTime, UpdateAt: sampleTime}
	mock.ExpectQuery("SELECT id, username, email, created_at, updated_at FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow(1, user.Username, user.Email, user.CreatedAt, user.UpdateAt))

	var p = persistence.NewPersistence(db)
	r, err := p.GetUserRepository().GetUsers()

	if err != nil {
		t.Errorf("error was not expected while getting users: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, u := range r {
		if u.Username != user.Username {
			t.Errorf("expected: %s, got: %s", user.Username, u.Username)
		}
	}
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sampleTime, err := time.Parse(time.DateOnly, "2021-01-01")

	user := User{ID: 1, Username: "user", Email: "mail", CreatedAt: sampleTime, UpdateAt: sampleTime}
	mock.ExpectQuery("SELECT id, username, email, created_at, updated_at FROM users").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow(1, user.Username, user.Email, user.CreatedAt, user.UpdateAt))

	var p = persistence.NewPersistence(db)
	r, err := p.GetUserRepository().GetUser(user.ID)

	if err != nil {
		t.Errorf("error was not expected while getting user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if r.Username != user.Username {
		t.Errorf("expected: %s, got: %s", user.Username, r.Username)
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sampleTime, err := time.Parse(time.DateOnly, "2021-01-01")

	user := User{ID: 1, Username: "user", Email: "mail", CreatedAt: sampleTime, UpdateAt: sampleTime}
	mock.ExpectExec("UPDATE users").
		WithArgs(user.Username, user.Email, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = persistence.NewPersistence(db)
	err = p.GetUserRepository().UpdateUser(user)

	if err != nil {
		t.Errorf("error was not expected while updating user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sampleTime, err := time.Parse(time.DateOnly, "2021-01-01")

	user := User{Username: "user", Email: "mail", CreatedAt: sampleTime, UpdateAt: sampleTime}
	mock.ExpectQuery("SELECT id, username, email, created_at, updated_at FROM users").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow(1, user.Username, user.Email, user.CreatedAt, user.UpdateAt))

	var p = persistence.NewPersistence(db)
	r, err := p.GetUserRepository().GetUserByEmail(user.Email)

	if err != nil {
		t.Errorf("error was not expected while getting user by email: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if r.Username != user.Username {
		t.Errorf("expected: %s, got: %s", user.Username, r.Username)
	}
}
