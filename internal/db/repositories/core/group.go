package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/prulloac/fineasy/pkg"
)

type Group struct {
	pkg.Model
	Name      string `json:"name" validate:"required,min=1"`
	CreatedBy uint   `json:"created_by" validate:"required,min=1"`
}

func (g *Group) String() string {
	out, err := json.Marshal(g)
	if err != nil {
		return fmt.Sprintf("%+v", g.ID)
	}
	return string(out)
}

func (s *CoreRepository) CreateGroup(name string, createdBy uint) (*Group, error) {
	var g Group
	err := s.Persistence.QueryRow(`
	INSERT INTO groups (name, created_by)
	VALUES ($1, $2)
	RETURNING id, name, created_by, created_at, updated_at
	`, name, createdBy).Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_, err = s.InsertUserGroup(createdBy, g.ID, pkg.Accepted)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *CoreRepository) GetGroupByID(groupID uint) (*Group, error) {
	var g Group
	err := s.Persistence.QueryRow(`
	SELECT id, name, created_by, created_at, updated_at
	FROM groups
	WHERE id = $1
	`, groupID).Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *CoreRepository) GetGroupByUserID(gid, uid uint) (*Group, error) {
	var g Group
	err := s.Persistence.QueryRow(`
	SELECT g.id, g.name, g.created_by, g.created_at, g.updated_at
	FROM groups g
	JOIN user_groups ug ON g.id = ug.group_id
	WHERE g.id = $1 AND ug.user_id = $2
	`, gid, uid).Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *CoreRepository) GetGroupsByUserID(userID uint) ([]Group, error) {
	rows, err := s.Persistence.Query(`
	SELECT g.id, g.name, g.created_by, g.created_at, g.updated_at
	FROM groups g
	JOIN user_groups ug ON g.id = ug.group_id
	WHERE ug.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []Group{}
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (s *CoreRepository) UpdateGroup(groupID uint, name string) (*Group, error) {
	var g Group
	err := s.Persistence.QueryRow(`
	UPDATE groups
	SET name = $2
	WHERE id = $1
	RETURNING id, name, created_by, created_at, updated_at
	`, groupID, name).Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
