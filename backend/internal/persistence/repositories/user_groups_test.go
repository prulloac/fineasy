package repositories

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/prulloac/fineasy/internal/persistence/entity"
)

func TestInsertUserGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	userGroup := entity.UserGroup{UserID: 1, GroupID: 1}
	mock.ExpectQuery("SELECT id FROM user_groups").
		WithArgs(userGroup.UserID, userGroup.GroupID).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO user_groups").
		WithArgs(userGroup.UserID, userGroup.GroupID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = UserGroupsRepository{db}
	err = p.Insert(userGroup)

	if err != nil {
		t.Errorf("error was not expected while inserting user group: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllUserGroups(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	userGroup := entity.UserGroup{ID: 1, UserID: 1, GroupID: 1}
	mock.ExpectQuery("SELECT id, user_id, group_id FROM user_groups").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "group_id"}).
			AddRow(userGroup.ID, userGroup.UserID, userGroup.GroupID))

	var p = UserGroupsRepository{db}
	r, err := p.GetAll()

	if err != nil {
		t.Errorf("error was not expected while getting user groups: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, e := range r {
		if e.ID != userGroup.ID {
			t.Errorf("expected id %d, but got %d", userGroup.ID, e.ID)
		}
		if e.UserID != userGroup.UserID {
			t.Errorf("expected user id %d, but got %d", userGroup.UserID, e.UserID)
		}
		if e.GroupID != userGroup.GroupID {
			t.Errorf("expected group id %d, but got %d", userGroup.GroupID, e.GroupID)
		}
	}
}

func TestGetUserGroupsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	userGroup := entity.UserGroup{ID: 1, UserID: 1, GroupID: 1}
	mock.ExpectQuery("SELECT id, user_id, group_id FROM user_groups").
		WithArgs(userGroup.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "group_id"}).
			AddRow(userGroup.ID, userGroup.UserID, userGroup.GroupID))

	var p = UserGroupsRepository{db}
	r, err := p.GetByUserID(userGroup.UserID)

	if err != nil {
		t.Errorf("error was not expected while getting user groups by user id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, e := range r {
		if e.UserID != userGroup.UserID {
			t.Errorf("expected user id %d, but got %d", userGroup.UserID, e.UserID)
		}
	}
}

func TestGetUserGroupsByGroupID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	userGroup := entity.UserGroup{ID: 1, UserID: 1, GroupID: 1}
	mock.ExpectQuery("SELECT id, user_id, group_id FROM user_groups").
		WithArgs(userGroup.GroupID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "group_id"}).
			AddRow(userGroup.ID, userGroup.UserID, userGroup.GroupID))

	var p = UserGroupsRepository{db}
	r, err := p.GetByGroupID(userGroup.GroupID)

	if err != nil {
		t.Errorf("error was not expected while getting user groups by group id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, e := range r {
		if e.GroupID != userGroup.GroupID {
			t.Errorf("expected group id %d, but got %d", userGroup.GroupID, e.GroupID)
		}
	}
}

func TestGetUserGroupByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	userGroup := entity.UserGroup{ID: 1, UserID: 1, GroupID: 1}
	mock.ExpectQuery("SELECT id, user_id, group_id FROM user_groups").
		WithArgs(userGroup.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "group_id"}).
			AddRow(userGroup.ID, userGroup.UserID, userGroup.GroupID))

	var p = UserGroupsRepository{db}
	r, err := p.GetByID(userGroup.ID)

	if err != nil {
		t.Errorf("error was not expected while getting user group by id: %s", err)
	}

	if r.ID != userGroup.ID {
		t.Errorf("expected id %d, but got %d", userGroup.ID, r.ID)
	}
	if r.UserID != userGroup.UserID {
		t.Errorf("expected user id %d, but got %d", userGroup.UserID, r.UserID)
	}
	if r.GroupID != userGroup.GroupID {
		t.Errorf("expected group id %d, but got %d", userGroup.GroupID, r.GroupID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
