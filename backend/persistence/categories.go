package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Category struct {
	ID          int
	Name        string
	Icon        string
	Color       string
	Description string
	DeletedAt   sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
	err := p.db.QueryRow("SELECT id FROM categories WHERE name = $1", category.Name).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := p.db.Exec("INSERT INTO categories (name, icon, color, description) VALUES ($1, $2, $3, $4)", category.Name, category.Icon, category.Color, category.Description)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (p *Persistence) GetCategories() ([]Category, error) {
	rows, err := p.db.Query("SELECT id, name, icon, color, description FROM categories WHERE deleted_at is NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name,
			&category.Icon, &category.Color, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (p *Persistence) GetCategory(id int) (Category, error) {
	var category Category
	err := p.db.QueryRow("SELECT id, name, icon, color, description, deleted_at, created_at, updated_at FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name,
		&category.Icon, &category.Color, &category.Description, &category.DeletedAt, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (p *Persistence) UpdateCategory(category Category) error {
	_, err := p.db.Exec("UPDATE categories SET name = $1, icon = $2, color = $3, description = $4 WHERE id = $5", category.Name, category.Icon, category.Color, category.Description, category.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Persistence) DeleteCategory(id int) error {
	// use soft delete
	_, err := p.db.Exec("UPDATE categories SET deleted_at = NOW() WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (p *Persistence) RestoreCategory(id int) (Category, error) {
	var c Category
	_, err := p.db.Exec("UPDATE categories SET deleted_at = NULL WHERE id = $1", id)
	if err != nil {
		return c, err
	}
	c, err = p.GetCategory(id)
	if err != nil {
		return c, err
	}
	return c, nil
}
