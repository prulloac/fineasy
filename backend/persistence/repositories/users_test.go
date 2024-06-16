package repositories

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/prulloac/fineasy/persistence/entity"
)

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	user := User{Username: "user", Email: "mail", CreatedAt: time.Now(), UpdateAt: time.Now()}
	mock.ExpectQuery("SELECT id FROM users").
		WithArgs(user.Email).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Username, user.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = UserRepository{db}
	err = p.InsertUser(user)

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

	user := User{Username: "user", Email: "mail", CreatedAt: time.Now(), UpdateAt: time.Now()}
	mock.ExpectQuery("SELECT id, username, email, created_at, updated_at FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow(1, user.Username, user.Email, user.CreatedAt, user.UpdateAt))

	var p = UserRepository{db}
	r, err := p.GetUsers()

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

	user := User{ID: 1, Username: "user", Email: "mail", CreatedAt: time.Now(), UpdateAt: time.Now()}
	mock.ExpectQuery("SELECT id, username, email, created_at, updated_at FROM users").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow(1, user.Username, user.Email, user.CreatedAt, user.UpdateAt))

	var p = UserRepository{db}
	r, err := p.GetUser(user.ID)

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

	user := User{ID: 1, Username: "user", Email: "mail", CreatedAt: time.Now(), UpdateAt: time.Now()}
	mock.ExpectExec("UPDATE users").
		WithArgs(user.Username, user.Email, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = UserRepository{db}
	err = p.UpdateUser(user)

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

	user := User{Username: "user", Email: "mail", CreatedAt: time.Now(), UpdateAt: time.Now()}
	mock.ExpectQuery("SELECT id, username, email, created_at, updated_at FROM users").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
			AddRow(1, user.Username, user.Email, user.CreatedAt, user.UpdateAt))

	var p = UserRepository{db}
	r, err := p.GetUserByEmail(user.Email)

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
