package auth

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	user := User{Username: "user", Hash: "hash", Email: "test@email.com", CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id FROM users").
		WithArgs(user.Email).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Hash, user.Username, user.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.InsertUser(user)

	if err != nil {
		t.Errorf("error was not expected while inserting user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	user := User{ID: 1, Username: "user", Hash: "hash", Email: "test@email.com", CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty, Disabled: false, ValidatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, hash, username, email, validated_at, disabled, created_at, updated_at FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "hash", "username", "email", "validated_at", "disabled", "created_at", "updated_at"}).
			AddRow(1, user.Hash, user.Username, user.Email, user.ValidatedAt, user.Disabled, user.CreatedAt, user.UpdatedAt))

	var p = AuthRepository{db}
	r, err := p.GetAllUsers()

	if err != nil {
		t.Errorf("error was not expected while getting all users: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, u := range r {
		if !reflect.DeepEqual(u, user) {
			t.Errorf("expected %v but got %v", user, u)
		}
	}
}
