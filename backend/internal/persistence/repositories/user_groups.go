package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/prulloac/fineasy/internal/persistence/entity"
)

type UserGroupsRepository struct {
	db *sql.DB
}

func NewUserGroupsRepository(db *sql.DB) *UserGroupsRepository {
	return &UserGroupsRepository{db}
}

func (ug *UserGroupsRepository) CreateTable() {
	data, _ := os.ReadFile("internal/persistence/schema/user_groups.sql")

	if data == nil {
		panic("Error reading user_groups schema file!")
	}

	_, err := ug.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating user_groups table!")
		panic(err)
	}
	fmt.Println("UserGroups table created!")
}

func (ug *UserGroupsRepository) DropTable() {
	_, err := ug.db.Exec("DROP TABLE IF EXISTS user_groups")
	if err != nil {
		fmt.Println("Error dropping user_groups table!")
		panic(err)
	}
	fmt.Println("UserGroups table dropped!")
}

func (ug *UserGroupsRepository) Insert(userGroup entity.UserGroup) error {
	// check if the userGroup already exists
	var id int
	err := ug.db.QueryRow(`
	SELECT
		id
	FROM user_groups
	WHERE user_id = $1 AND group_id = $2`, userGroup.UserID, userGroup.GroupID).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := ug.db.Exec(`
		INSERT INTO user_groups
		(user_id, group_id) VALUES ($1, $2)`,
			userGroup.UserID, userGroup.GroupID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ug *UserGroupsRepository) GetAll() ([]entity.UserGroup, error) {
	rows, err := ug.db.Query(`
	SELECT
		id, user_id, group_id
	FROM user_groups`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userGroups []entity.UserGroup
	for rows.Next() {
		var userGroup entity.UserGroup
		err := rows.Scan(&userGroup.ID, &userGroup.UserID, &userGroup.GroupID)
		if err != nil {
			return nil, err
		}
		userGroups = append(userGroups, userGroup)
	}
	return userGroups, nil
}

func (ug *UserGroupsRepository) GetByUserID(userID int) ([]entity.UserGroup, error) {
	rows, err := ug.db.Query(`
	SELECT
		id, user_id, group_id
	FROM user_groups
	WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userGroups []entity.UserGroup
	for rows.Next() {
		var userGroup entity.UserGroup
		err := rows.Scan(&userGroup.ID, &userGroup.UserID, &userGroup.GroupID)
		if err != nil {
			return nil, err
		}
		userGroups = append(userGroups, userGroup)
	}
	return userGroups, nil
}

func (ug *UserGroupsRepository) GetByGroupID(groupID int) ([]entity.UserGroup, error) {
	rows, err := ug.db.Query(`
	SELECT
		id, user_id, group_id
	FROM user_groups
	WHERE group_id = $1`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userGroups []entity.UserGroup
	for rows.Next() {
		var userGroup entity.UserGroup
		err := rows.Scan(&userGroup.ID, &userGroup.UserID, &userGroup.GroupID)
		if err != nil {
			return nil, err
		}
		userGroups = append(userGroups, userGroup)
	}
	return userGroups, nil
}

func (ug *UserGroupsRepository) GetByID(id int) (entity.UserGroup, error) {
	var userGroup entity.UserGroup
	err := ug.db.QueryRow(`
	SELECT
		id, user_id, group_id
	FROM user_groups
	WHERE id = $1`, id).Scan(&userGroup.ID, &userGroup.UserID, &userGroup.GroupID)
	if err != nil {
		return entity.UserGroup{}, err
	}
	return userGroup, nil
}
