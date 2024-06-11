package main

import (
	"fmt"

	"github.com/prulloac/fineasy/persistence"
	"github.com/prulloac/fineasy/persistence/categories"
	"github.com/prulloac/fineasy/routes"
)

func main() {
	db := persistence.Connection()
	defer db.Close()
	persistence.VerifySchema(db)
	categories.InsertCategory(db, categories.Category{Name: "Food", Icon: "restaurant", Color: "red", Description: "Food and drinks"})
	_categories := categories.GetCategories(db)
	if len(_categories) == 0 {
		panic("No categories found")
	}
	for _, category := range _categories {
		fmt.Println(category)
	}
	fmt.Println(categories.GetCategory(db, 1))
	categories.UpdateCategory(db, categories.Category{ID: 1, Name: "Food and Drinks", Icon: "restaurant", Color: "red", Description: "Food and drinks"})
	fmt.Println(categories.GetCategory(db, 1))
	categories.DeleteCategory(db, 1)
	fmt.Println(categories.GetCategory(db, 1))
	categories.RestoreCategory(db, 1)
	fmt.Println(categories.GetCategory(db, 1))
	fmt.Println("Server is running...")
	routes.Run()
}
