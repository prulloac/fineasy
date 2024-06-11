package categories

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
	DeletedAt   sql.NullString
	CreatedAt   string
	UpdatedAt   string
}

func CreateCategoriesTable(db *sql.DB) {
	data, _ := os.ReadFile("persistence/schema/categories.sql")
	_, err := db.Exec(string(data))
	if err != nil {
		panic(err)
	}
	fmt.Println("Categories table created!")
}

func InsertCategory(db *sql.DB, category Category) {
	// check if the category already exists regardless of the icon, color, and description
	var id int
	err := db.QueryRow("SELECT id FROM categories WHERE name = $1", category.Name).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := db.Exec("INSERT INTO categories (name, icon, color, description) VALUES ($1, $2, $3, $4)", category.Name, category.Icon, category.Color, category.Description)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
}

func GetCategories(db *sql.DB) []Category {
	rows, err := db.Query("SELECT id, name, icon, color, description FROM categories WHERE deleted_at is NULL")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name,
			&category.Icon, &category.Color, &category.Description)
		if err != nil {
			panic(err)
		}
		categories = append(categories, category)
	}
	return categories
}

func GetCategory(db *sql.DB, id int) Category {
	var category Category
	err := db.QueryRow("SELECT * FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name,
		&category.Icon, &category.Color, &category.Description, &category.DeletedAt, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		panic(err)
	}
	return category
}

func UpdateCategory(db *sql.DB, category Category) {
	_, err := db.Exec("UPDATE categories SET name = $1, icon = $2, color = $3, description = $4 WHERE id = $5", category.Name, category.Icon, category.Color, category.Description, category.ID)
	if err != nil {
		panic(err)
	}
}

func DeleteCategory(db *sql.DB, id int) {
	// use soft delete
	_, err := db.Exec("UPDATE categories SET deleted_at = NOW() WHERE id = $1", id)
	if err != nil {
		panic(err)
	}
}

func RestoreCategory(db *sql.DB, id int) {
	_, err := db.Exec("UPDATE categories SET deleted_at = NULL WHERE id = $1", id)
	if err != nil {
		panic(err)
	}
}
