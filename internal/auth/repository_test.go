package auth

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/prulloac/fineasy/tests"
)

func TestCreateAndDropTables(t *testing.T) {
	ctx := context.Background()
	container, err := tests.StartPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	db, err := container.DB()
	if err != nil {
		t.Fatal(err)
	}

	var p = AuthRepository{db}

	err = p.CreateTable()
	if err != nil {
		t.Errorf("error was not expected while creating tables: %s", err)
	}

	err = p.DropTable()
	if err != nil {
		t.Errorf("error was not expected while dropping tables: %s", err)
	}
}

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	hash, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}

	user := User{Username: "user", Hash: hash.String(), Email: "test@email.com", CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
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
	hash, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}

	user := User{ID: 1, Username: "user", Hash: hash.String(), Email: "test@email.com", CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty, Disabled: false, ValidatedAt: twentyTwenty}
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

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}
	hash, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}

	user := User{ID: 1, Username: "user", Hash: hash.String(), Email: "test@email.com", CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty, Disabled: false, ValidatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, hash, username, email, validated_at, disabled, created_at, updated_at FROM users").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "hash", "username", "email", "validated_at", "disabled", "created_at", "updated_at"}).
			AddRow(user.ID, user.Hash, user.Username, user.Email, user.ValidatedAt, user.Disabled, user.CreatedAt, user.UpdatedAt))

	var p = AuthRepository{db}
	r, err := p.GetUserByID(user.ID)

	if err != nil {
		t.Errorf("error was not expected while getting user by id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !reflect.DeepEqual(r, user) {
		t.Errorf("expected %v but got %v", user, r)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}
	hash, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}

	user := User{ID: 1, Username: "user", Hash: hash.String(), Email: "test@email.com", CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty, Disabled: false, ValidatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, hash, username, email, validated_at, disabled, created_at, updated_at FROM users").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "hash", "username", "email", "validated_at", "disabled", "created_at", "updated_at"}).
			AddRow(user.ID, user.Hash, user.Username, user.Email, user.ValidatedAt, user.Disabled, user.CreatedAt, user.UpdatedAt))

	var p = AuthRepository{db}
	r, err := p.GetUserByEmail(user.Email)

	if err != nil {
		t.Errorf("error was not expected while getting user by email: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !reflect.DeepEqual(r, user) {
		t.Errorf("expected %v but got %v", user, r)
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}
	hash, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}

	user := User{ID: 1, Username: "user", Hash: hash.String(), Email: "test@email.com", CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty, Disabled: false, ValidatedAt: twentyTwenty}
	mock.ExpectExec("UPDATE users").
		WithArgs(user.Username, user.ValidatedAt, user.Disabled, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.UpdateUser(user)

	if err != nil {
		t.Errorf("error was not expected while updating user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	userID := 1
	mock.ExpectExec("DELETE FROM users").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.DeleteUser(userID)

	if err != nil {
		t.Errorf("error was not expected while deleting user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertInternalLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	hash, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}
	salt, err := uuid.NewV7FromReader(strings.NewReader(hash.String()))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}
	if salt == hash {
		t.Fatalf("expected different uuids but got the same")
	}

	login := InternalLogin{UserID: 1, Email: "email", Password: "password", PasswordSalt: salt.String(), Algorithm: "algorithm", PasswordLastUpdatedAt: twentyTwenty, LoginAttempts: 1, LastLoginAttempt: twentyTwenty, LastLoginSuccess: twentyTwenty}
	mock.ExpectQuery("SELECT id FROM internal_logins").
		WithArgs(login.Email).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO internal_logins").
		WithArgs(login.UserID, login.Email, login.Password, login.PasswordSalt, login.Algorithm).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.InsertInternalLogin(login)

	if err != nil {
		t.Errorf("error was not expected while inserting internal login: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllInternalLogins(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	hash, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}
	salt, err := uuid.NewV7FromReader(strings.NewReader(hash.String()))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}
	if salt == hash {
		t.Fatalf("expected different uuids but got the same")
	}

	login := InternalLogin{ID: 1, UserID: 1, Email: "test@email.com", Password: "password", PasswordSalt: salt.String(), Algorithm: "algorithm", PasswordLastUpdatedAt: twentyTwenty, LoginAttempts: 1, LastLoginAttempt: twentyTwenty, LastLoginSuccess: twentyTwenty, CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, email, password, password_salt, algorithm, password_last_updated_at, login_attempts, last_login_attempt, last_login_success, created_at, updated_at FROM internal_logins").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "email", "password", "password_salt", "algorithm", "password_last_updated_at", "login_attempts", "last_login_attempt", "last_login_success", "created_at", "updated_at"}).
			AddRow(login.ID, login.UserID, login.Email, login.Password, login.PasswordSalt, login.Algorithm, login.PasswordLastUpdatedAt, login.LoginAttempts, login.LastLoginAttempt, login.LastLoginSuccess, twentyTwenty, twentyTwenty))

	var p = AuthRepository{db}
	r, err := p.GetAllInternalLogins()

	if err != nil {
		t.Errorf("error was not expected while getting all internal logins: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, l := range r {
		if !reflect.DeepEqual(l, login) {
			t.Errorf("expected %v but got %v", login, l)
		}
	}
}

func TestGetInternalLoginByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	hash, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}
	salt, err := uuid.NewV7FromReader(strings.NewReader(hash.String()))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}
	if salt == hash {
		t.Fatalf("expected different uuids but got the same")
	}

	login := InternalLogin{ID: 1, UserID: 1, Email: "test@email.com", Password: "password", PasswordSalt: salt.String(), Algorithm: "algorithm", PasswordLastUpdatedAt: twentyTwenty, LoginAttempts: 1, LastLoginAttempt: twentyTwenty, LastLoginSuccess: twentyTwenty, CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, email, password, password_salt, algorithm, password_last_updated_at, login_attempts, last_login_attempt, last_login_success, created_at, updated_at FROM internal_logins").
		WithArgs(login.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "email", "password", "password_salt", "algorithm", "password_last_updated_at", "login_attempts", "last_login_attempt", "last_login_success", "created_at", "updated_at"}).
			AddRow(login.ID, login.UserID, login.Email, login.Password, login.PasswordSalt, login.Algorithm, login.PasswordLastUpdatedAt, login.LoginAttempts, login.LastLoginAttempt, login.LastLoginSuccess, twentyTwenty, twentyTwenty))

	var p = AuthRepository{db}
	r, err := p.GetInternalLoginByID(login.ID)

	if err != nil {
		t.Errorf("error was not expected while getting internal login by id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !reflect.DeepEqual(r, login) {
		t.Errorf("expected %v but got %v", login, r)
	}
}

func TestGetInternalLoginByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	hash, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}
	salt, err := uuid.NewV7FromReader(strings.NewReader(hash.String()))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}
	if salt == hash {
		t.Fatalf("expected different uuids but got the same")
	}

	login := InternalLogin{ID: 1, UserID: 1, Email: "test@email.com", Password: "password", PasswordSalt: salt.String(), Algorithm: "algorithm", PasswordLastUpdatedAt: twentyTwenty, LoginAttempts: 1, LastLoginAttempt: twentyTwenty, LastLoginSuccess: twentyTwenty, CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, email, password, password_salt, algorithm, password_last_updated_at, login_attempts, last_login_attempt, last_login_success, created_at, updated_at FROM internal_logins").
		WithArgs(login.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "email", "password", "password_salt", "algorithm", "password_last_updated_at", "login_attempts", "last_login_attempt", "last_login_success", "created_at", "updated_at"}).
			AddRow(login.ID, login.UserID, login.Email, login.Password, login.PasswordSalt, login.Algorithm, login.PasswordLastUpdatedAt, login.LoginAttempts, login.LastLoginAttempt, login.LastLoginSuccess, twentyTwenty, twentyTwenty))

	var p = AuthRepository{db}
	r, err := p.GetInternalLoginByEmail(login.Email)

	if err != nil {
		t.Errorf("error was not expected while getting internal login by email: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !reflect.DeepEqual(r, login) {
		t.Errorf("expected %v but got %v", login, r)
	}
}

func TestUpdateInternalLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	hash, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}
	salt, err := uuid.NewV7FromReader(strings.NewReader(hash.String()))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating uuid", err)
	}
	if salt == hash {
		t.Fatalf("expected different uuids but got the same")
	}

	login := InternalLogin{ID: 1, UserID: 1, Email: "test@email.com", Password: "password", PasswordSalt: salt.String(), Algorithm: "algorithm", PasswordLastUpdatedAt: twentyTwenty, LoginAttempts: 1, LastLoginAttempt: twentyTwenty, LastLoginSuccess: twentyTwenty, CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectExec("UPDATE internal_logins").
		WithArgs(login.Password, login.Algorithm, login.PasswordLastUpdatedAt, login.LoginAttempts, login.LastLoginAttempt, login.LastLoginSuccess, login.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.UpdateInternalLogin(login)

	if err != nil {
		t.Errorf("error was not expected while updating internal login: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteInternalLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	loginID := 1
	mock.ExpectExec("DELETE FROM internal_logins").
		WithArgs(loginID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.DeleteInternalLogin(loginID)

	if err != nil {
		t.Errorf("error was not expected while deleting internal login: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertLoginToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	token := LoginToken{UserID: 1, Token: "token", TokenType: 1, ExpiresAt: twentyTwenty, UsedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id FROM login_tokens").
		WithArgs(token.Token).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO login_tokens").
		WithArgs(token.UserID, token.Token, token.TokenType, token.ExpiresAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.InsertLoginToken(token)

	if err != nil {
		t.Errorf("error was not expected while inserting login token: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllLoginTokens(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	token := LoginToken{ID: 1, UserID: 1, Token: "token", TokenType: 1, ExpiresAt: twentyTwenty, UsedAt: twentyTwenty, CreatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, token, token_type, expires_at, used_at, created_at FROM login_tokens").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "token_type", "expires_at", "used_at", "created_at"}).
			AddRow(token.ID, token.UserID, token.Token, token.TokenType, token.ExpiresAt, token.UsedAt, token.CreatedAt))

	var p = AuthRepository{db}
	r, err := p.GetAllLoginTokens()

	if err != nil {
		t.Errorf("error was not expected while getting all login tokens: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, to := range r {
		if !reflect.DeepEqual(to, token) {
			t.Errorf("expected %v but got %v", token, to)
		}
	}
}

func TestGetAllLoginTokensByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	userID := 1
	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	token := LoginToken{ID: 1, UserID: 1, Token: "token", TokenType: 1, ExpiresAt: twentyTwenty, UsedAt: twentyTwenty, CreatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, token, token_type, expires_at, used_at, created_at FROM login_tokens").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "token_type", "expires_at", "used_at", "created_at"}).
			AddRow(token.ID, token.UserID, token.Token, token.TokenType, token.ExpiresAt, token.UsedAt, token.CreatedAt))

	var p = AuthRepository{db}
	r, err := p.GetAllLoginTokensByUserID(userID)

	if err != nil {
		t.Errorf("error was not expected while getting all login tokens by user id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, to := range r {
		if !reflect.DeepEqual(to, token) {
			t.Errorf("expected %v but got %v", token, to)
		}
	}
}

func TestGetLoginTokenByTokenAndUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	userID := 1
	token := "token"
	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	loginToken := LoginToken{ID: 1, UserID: 1, Token: "token", TokenType: 1, ExpiresAt: twentyTwenty, UsedAt: twentyTwenty, CreatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, token, token_type, expires_at, used_at, created_at FROM login_tokens").
		WithArgs(token, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "token_type", "expires_at", "used_at", "created_at"}).
			AddRow(loginToken.ID, loginToken.UserID, loginToken.Token, loginToken.TokenType, loginToken.ExpiresAt, loginToken.UsedAt, loginToken.CreatedAt))

	var p = AuthRepository{db}
	r, err := p.GetLoginTokenByTokenAndUserID(token, userID)

	if err != nil {
		t.Errorf("error was not expected while getting login token by token and user id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !reflect.DeepEqual(r, loginToken) {
		t.Errorf("expected %v but got %v", loginToken, r)
	}
}

func TestUpdateLoginToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	token := LoginToken{ID: 1, UserID: 1, Token: "token", TokenType: 1, ExpiresAt: twentyTwenty, UsedAt: twentyTwenty, CreatedAt: twentyTwenty}
	mock.ExpectExec("UPDATE login_tokens").
		WithArgs(token.Token, token.TokenType, token.ExpiresAt, token.UsedAt, token.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.UpdateLoginToken(token)

	if err != nil {
		t.Errorf("error was not expected while updating login token: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteLoginToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	tokenID := 1
	mock.ExpectExec("DELETE FROM login_tokens").
		WithArgs(tokenID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.DeleteLoginToken(tokenID)

	if err != nil {
		t.Errorf("error was not expected while deleting login token: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertExternalLoginProvider(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	provider := ExternalLoginProvider{ID: 1, Name: "name", Type: 1, Endpoint: "endpoint", Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mock.ExpectQuery("SELECT id FROM external_login_providers").
		WithArgs(provider.Name, provider.Type, provider.Endpoint).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO external_login_providers").
		WithArgs(provider.Name, provider.Type, provider.Endpoint).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.InsertExternalLoginProvider(provider)

	if err != nil {
		t.Errorf("error was not expected while inserting external login provider: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllExternalLoginProviders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	provider := ExternalLoginProvider{ID: 1, Name: "name", Type: 1, Endpoint: "http://localhost:8080", Enabled: true, CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, name, type, endpoint, enabled, created_at, updated_at FROM external_login_providers").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "endpoint", "enabled", "created_at", "updated_at"}).
			AddRow(provider.ID, provider.Name, provider.Type, provider.Endpoint, provider.Enabled, provider.CreatedAt, provider.UpdatedAt))

	var p = AuthRepository{db}
	r, err := p.GetAllExternalLoginProviders()

	if err != nil {
		t.Errorf("error was not expected while getting all external login providers: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, pr := range r {
		if !reflect.DeepEqual(pr, provider) {
			t.Errorf("expected %v but got %v", provider, pr)
		}
	}
}

func TestGetAllExternalLoginProvidersByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	name := "name"
	provider := ExternalLoginProvider{ID: 1, Name: "name", Type: 1, Endpoint: "http://localhost:8080", Enabled: true, CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, name, type, endpoint, enabled, created_at, updated_at FROM external_login_providers").
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "endpoint", "enabled", "created_at", "updated_at"}).
			AddRow(provider.ID, provider.Name, provider.Type, provider.Endpoint, provider.Enabled, provider.CreatedAt, provider.UpdatedAt))

	var p = AuthRepository{db}
	r, err := p.GetAllExternalLoginProvidersByName(name)

	if err != nil {
		t.Errorf("error was not expected while getting all external login providers by name: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, pr := range r {
		if !reflect.DeepEqual(pr, provider) {
			t.Errorf("expected %v but got %v", provider, pr)
		}
	}
}

func TestGetExternalLoginProviderByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	providerID := 1
	provider := ExternalLoginProvider{ID: 1, Name: "name", Type: 1, Endpoint: "http://localhost:8080", Enabled: true, CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, name, type, endpoint, enabled, created_at, updated_at FROM external_login_providers").
		WithArgs(providerID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "endpoint", "enabled", "created_at", "updated_at"}).
			AddRow(provider.ID, provider.Name, provider.Type, provider.Endpoint, provider.Enabled, provider.CreatedAt, provider.UpdatedAt))

	var p = AuthRepository{db}
	r, err := p.GetExternalLoginProviderByID(providerID)

	if err != nil {
		t.Errorf("error was not expected while getting external login provider by id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !reflect.DeepEqual(r, provider) {
		t.Errorf("expected %v but got %v", provider, r)
	}
}

func TestUpdateExternalLoginProvider(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	provider := ExternalLoginProvider{ID: 1, Name: "name", Type: 1, Endpoint: "http://localhost:8080", Enabled: true, CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectExec("UPDATE external_login_providers").
		WithArgs(provider.Name, provider.Type, provider.Endpoint, provider.Enabled, provider.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.UpdateExternalLoginProvider(provider)

	if err != nil {
		t.Errorf("error was not expected while updating external login provider: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteExternalLoginProviderByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	providerID := 1
	mock.ExpectExec("DELETE FROM external_login_providers").
		WithArgs(providerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.DeleteExternalLoginProvider(providerID)

	if err != nil {
		t.Errorf("error was not expected while deleting external login provider by id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertExternalLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	login := ExternalLogin{ID: 1, UserID: 1, ProviderID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mock.ExpectQuery("SELECT id FROM external_logins").
		WithArgs(login.UserID, login.ProviderID).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO external_logins").
		WithArgs(login.UserID, login.ProviderID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.InsertExternalLogin(login)

	if err != nil {
		t.Errorf("error was not expected while inserting external login: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllExternalLogins(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	login := ExternalLogin{ID: 1, UserID: 1, ProviderID: 1, CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, provider_id, created_at, updated_at FROM external_logins").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "provider_id", "created_at", "updated_at"}).
			AddRow(login.ID, login.UserID, login.ProviderID, login.CreatedAt, login.UpdatedAt))

	var p = AuthRepository{db}
	r, err := p.GetAllExternalLogins()

	if err != nil {
		t.Errorf("error was not expected while getting all external logins: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, l := range r {
		if !reflect.DeepEqual(l, login) {
			t.Errorf("expected %v but got %v", login, l)
		}
	}
}

func TestGetAllExternalLoginsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	userID := 1
	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	login := ExternalLogin{ID: 1, UserID: 1, ProviderID: 1, CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, provider_id, created_at, updated_at FROM external_logins").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "provider_id", "created_at", "updated_at"}).
			AddRow(login.ID, login.UserID, login.ProviderID, login.CreatedAt, login.UpdatedAt))

	var p = AuthRepository{db}
	r, err := p.GetAllExternalLoginsByUserID(userID)

	if err != nil {
		t.Errorf("error was not expected while getting all external logins by user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, l := range r {
		if !reflect.DeepEqual(l, login) {
			t.Errorf("expected %v but got %v", login, l)
		}
	}
}

func TestGetExternalLoginByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	loginID := 1

	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	login := ExternalLogin{ID: 1, UserID: 1, ProviderID: 1, CreatedAt: twentyTwenty, UpdatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, provider_id, created_at, updated_at FROM external_logins").
		WithArgs(loginID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "provider_id", "created_at", "updated_at"}).
			AddRow(login.ID, login.UserID, login.ProviderID, login.CreatedAt, login.UpdatedAt))

	var p = AuthRepository{db}
	r, err := p.GetExternalLoginByID(loginID)

	if err != nil {
		t.Errorf("error was not expected while getting external login by id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !reflect.DeepEqual(r, login) {
		t.Errorf("expected %v but got %v", login, r)
	}
}

func TestDeleteExternalLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	loginID := 1
	mock.ExpectExec("DELETE FROM external_logins").
		WithArgs(loginID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = AuthRepository{db}
	err = p.DeleteExternalLogin(loginID)

	if err != nil {
		t.Errorf("error was not expected while deleting external login by id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertExternalLoginToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}
	token := ExternalLoginToken{ID: 1, UserID: 1, ProviderID: 1, Token: "token", LoginIP: "127.0.0.1", UserAgent: "user agent", LoggedInAt: twentyTwenty, CreatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id FROM external_login_tokens").
		WithArgs(token.Token).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO external_login_tokens").
		WithArgs(token.UserID, token.ProviderID, token.LoginIP, token.UserAgent, token.LoggedInAt, token.Token).
		WillReturnResult(sqlmock.NewResult(1, 1))
	var p = AuthRepository{db}
	err = p.InsertExternalLoginToken(token)
	if err != nil {
		t.Errorf("error was not expected while inserting external login token: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllExternalLoginTokens(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)

	}
	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)

	}
	token := ExternalLoginToken{ID: 1, UserID: 1, ProviderID: 1, Token: "token", LoginIP: "127.0.0.1", UserAgent: "user agent", LoggedInAt: twentyTwenty, CreatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, provider_id, login_ip, user_agent, logged_in_at, token, created_at FROM external_login_tokens").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "provider_id", "login_ip", "user_agent", "logged_in_at", "token", "created_at"}).
			AddRow(token.ID, token.UserID, token.ProviderID, token.LoginIP, token.UserAgent, token.LoggedInAt, token.Token, token.CreatedAt))

	var p = AuthRepository{db}
	r, err := p.GetAllExternalLoginTokens()
	if err != nil {
		t.Errorf("error was not expected while getting all external login tokens: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)

	}
	for _, to := range r {
		if !reflect.DeepEqual(to, token) {
			t.Errorf("expected %v but got %v", token, to)
		}
	}
}

func TestGetAllExternalLoginTokensByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)

	}
	userID := 1
	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)

	}
	token := ExternalLoginToken{ID: 1, UserID: 1, ProviderID: 1, Token: "token", LoginIP: "127.0.0.1", UserAgent: "user agent", LoggedInAt: twentyTwenty, CreatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, provider_id, login_ip, user_agent, logged_in_at, token, created_at FROM external_login_tokens").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "provider_id", "login_ip", "user_agent", "logged_in_at", "token", "created_at"}).
			AddRow(token.ID, token.UserID, token.ProviderID, token.LoginIP, token.UserAgent, token.LoggedInAt, token.Token, token.CreatedAt))

	var p = AuthRepository{db}
	r, err := p.GetAllExternalLoginTokensByUserID(userID)
	if err != nil {
		t.Errorf("error was not expected while getting all external login tokens by user id: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)

	}
	for _, to := range r {
		if !reflect.DeepEqual(to, token) {
			t.Errorf("expected %v but got %v", token, to)
		}
	}
}

func TestGetExternalLoginTokenByTokenAndUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)

	}
	userID := 1
	token := "token"
	twentyTwenty, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)

	}
	loginToken := ExternalLoginToken{ID: 1, UserID: 1, ProviderID: 1, Token: "token", LoginIP: "127.0.0.1", UserAgent: "user agent", LoggedInAt: twentyTwenty, CreatedAt: twentyTwenty}
	mock.ExpectQuery("SELECT id, user_id, provider_id, login_ip, user_agent, logged_in_at, token, created_at FROM external_login_tokens").
		WithArgs(token, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "provider_id", "login_ip", "user_agent", "logged_in_at", "token", "created_at"}).
			AddRow(loginToken.ID, loginToken.UserID, loginToken.ProviderID, loginToken.LoginIP, loginToken.UserAgent, loginToken.LoggedInAt, loginToken.Token, loginToken.CreatedAt))

	var p = AuthRepository{db}
	r, err := p.GetExternalLoginTokenByTokenAndUserID(token, userID)
	if err != nil {
		t.Errorf("error was not expected while getting external login token by token and user id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)

	}

	if !reflect.DeepEqual(r, loginToken) {
		t.Errorf("expected %v but got %v", loginToken, r)
	}

}

func TestDeleteExternalLoginToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)

	}
	tokenID := 1
	mock.ExpectExec("DELETE FROM external_login_tokens").
		WithArgs(tokenID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	var p = AuthRepository{db}
	err = p.DeleteExternalLoginToken(tokenID)
	if err != nil {
		t.Errorf("error was not expected while deleting external login token: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)

	}
}
