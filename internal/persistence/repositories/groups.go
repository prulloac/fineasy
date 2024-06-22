package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/prulloac/fineasy/internal/persistence/entity"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db}
}

func (g *GroupRepository) CreateTable() {
	data, _ := os.ReadFile("internal/persistence/schema/groups.sql")

	if data == nil {
		panic("Error reading accounts schema file!")
	}

	_, err := g.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating groups table!")
		panic(err)
	}
	fmt.Println("Categories table created!")
}

func (g *GroupRepository) DropTable() {
	_, err := g.db.Exec("DROP TABLE IF EXISTS groups")
	if err != nil {
		fmt.Println("Error dropping groups table!")
		panic(err)
	}
	fmt.Println("Groups table dropped!")
}

func (g *GroupRepository) Insert(group entity.Group) error {
	// check if the group already exists
	var id int
	err := g.db.QueryRow(`
	SELECT
		id
	FROM groups
	WHERE name = $1`, group.Name).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := g.db.Exec(`
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

func (g *GroupRepository) GetByUserID(userID int) ([]entity.Group, error) {
	rows, err := g.db.Query(`
	SELECT 
		id,
		name,
		created_by,
		created_at
	FROM groups
	WHERE created_by = $1`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []entity.Group
	for rows.Next() {
		var group entity.Group
		err := rows.Scan(&group.ID, &group.Name, &group.CreatedBy, &group.CreatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (g *GroupRepository) GetByID(id int) (entity.Group, error) {
	var group entity.Group
	err := g.db.QueryRow(`
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

func (g *GroupRepository) Update(group entity.Group) error {
	_, err := g.db.Exec(`
	UPDATE groups
	SET name = $1
	WHERE id = $2`, group.Name, group.ID)
	if err != nil {
		return err
	}
	return nil
}
