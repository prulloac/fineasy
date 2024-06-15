package persistence

import (
	"database/sql"
	"fmt"
	"os"
)

type Category struct {
	ID          int
	Name        string
	Icon        string
	Color       string
	Description string
	Order       int
	GroupID     int
}

func (p *Persistence) CreateCategoriesTable() {
	data, _ := os.ReadFile("persistence/schema/categories.sql")
	_, err := p.db.Exec(string(data))
	if err != nil {
		panic(err)
	}
	fmt.Println("Categories table created!")
}

func (p *Persistence) InsertCategory(category Category) error {
	// check if the category already exists regardless of the icon, color, and description
	var id int
	err := p.db.QueryRow(`
	SELECT 
		id 
	FROM categories 
	WHERE name = $1`, category.Name).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := p.db.Exec(`
		INSERT INTO categories 
		(name, icon, color, description, ord, group_id) VALUES ($1, $2, $3, $4, $5, $6)`,
			category.Name, category.Icon, category.Color,
			category.Description, category.Order, category.GroupID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (p *Persistence) GetCategories(group_id int) ([]Category, error) {
	rows, err := p.db.Query(`
	SELECT 
		id, 
		name, 
		icon, 
		color, 
		description, 
		ord,
		group_id
	FROM categories
	WHERE group_id = $1
	ORDER BY ord ASC
	`, group_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name,
			&category.Icon, &category.Color, &category.Description,
			&category.Order, &category.GroupID)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (p *Persistence) GetCategory(id int) (Category, error) {
	var category Category
	err := p.db.QueryRow(`
	SELECT 
		id, 
		name, 
		icon, 
		color, 
		description, 
		ord, 
		group_id 
	FROM categories 
	WHERE id = $1
	`, id).Scan(&category.ID, &category.Name, &category.Icon,
		&category.Color, &category.Description, &category.Order, &category.GroupID)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (p *Persistence) UpdateCategory(category Category) error {
	_, err := p.db.Exec(`
	UPDATE categories 
	SET 
		name = $1, 
		icon = $2, 
		color = $3, 
		description = $4, 
		ord = $5 
	WHERE id = $6`,
		category.Name, category.Icon, category.Color, category.Description, category.Order, category.ID)
	if err != nil {
		return err
	}
	return nil
}
