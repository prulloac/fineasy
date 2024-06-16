package repositories

import (
	"database/sql"
	"fmt"
	"os"

	. "github.com/prulloac/fineasy/persistence/entity"
)

type GroupRepository struct {
	DB *sql.DB
}

func (g *GroupRepository) CreateGroupsTable() {
	data, _ := os.ReadFile("persistence/schema/groups.sql")
	_, err := g.DB.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating groups table!")
		panic(err)
	}
	fmt.Println("Categories table created!")
}

func (g *GroupRepository) InsertGroup(group Group) error {
	// check if the group already exists
	var id int
	err := g.DB.QueryRow(`
	SELECT
		id
	FROM groups
	WHERE name = $1`, group.Name).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := g.DB.Exec(`
		INSERT INTO groups
		(name, created_by,) VALUES ($1, $2)`,
			group.Name, group.CreatedBy)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (g *GroupRepository) GetGroups(user_id int) ([]Group, error) {
	rows, err := g.DB.Query(`
	SELECT 
		id,
		name,
		created_by,
		created_at
	FROM groups
	WHERE created_by = $1`, user_id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.ID, &group.Name, &group.CreatedBy, &group.CreatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (g *GroupRepository) GetGroup(id int) (Group, error) {
	var group Group
	err := g.DB.QueryRow(`
	SELECT 
		id,
		name,
		created_by,
		created_at
	FROM groups
	WHERE id = $1`, id).Scan(&group.ID, &group.Name, &group.CreatedBy, &group.CreatedAt)

	if err != nil {
		return group, err
	}
	return group, nil
}

func (g *GroupRepository) UpdateGroup(group Group) error {
	_, err := g.DB.Exec(`
	UPDATE groups
	SET name = $1
	WHERE id = $2`, group.Name, group.ID)
	if err != nil {
		return err
	}
	return nil
}
