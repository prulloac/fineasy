package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/prulloac/fineasy/persistence/entity"
)

type CategoriesRepository struct {
	db *sql.DB
}

func NewCategoriesRepository(db *sql.DB) *CategoriesRepository {
	return &CategoriesRepository{db}
}

func (c *CategoriesRepository) CreateCategoriesTable() {
	data, _ := os.ReadFile("persistence/schema/categories.sql")
	_, err := c.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating categories table!")
		panic(err)
	}
	fmt.Println("Categories table created!")
}

func (c *CategoriesRepository) InsertCategory(category entity.Category) error {
	// check if the category already exists regardless of the icon, color, and description
	var id int
	err := c.db.QueryRow(`
	SELECT 
		id 
	FROM categories 
	WHERE name = $1`, category.Name).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := c.db.Exec(`
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

func (c *CategoriesRepository) GetCategories(groupID int) ([]entity.Category, error) {
	rows, err := c.db.Query(`
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
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []entity.Category
	for rows.Next() {
		var category entity.Category
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

func (c *CategoriesRepository) GetCategory(id int) (entity.Category, error) {
	var category entity.Category
	err := c.db.QueryRow(`
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

func (c *CategoriesRepository) UpdateCategory(category entity.Category) error {
	_, err := c.db.Exec(`
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
