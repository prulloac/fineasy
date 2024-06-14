package main

import (
	"fmt"

	"github.com/prulloac/fineasy/persistence"
	"github.com/prulloac/fineasy/routes"
)

func main() {
	p := persistence.Connect()
	defer p.Close()
	p.VerifySchema()
	p.InsertCategory(persistence.Category{Name: "Food", Icon: "restaurant", Color: "red", Description: "Food and drinks"})
	_categories, err := p.GetCategories()
	if err != nil {
		panic(err)
	}
	if len(_categories) == 0 {
		panic("No categories found")
	}
	for _, category := range _categories {
		fmt.Println(category)
	}
	fmt.Println(p.GetCategory(1))
	p.UpdateCategory(persistence.Category{ID: 1, Name: "Food and Drinks", Icon: "restaurant", Color: "red", Description: "Food and drinks"})
	fmt.Println(p.GetCategory(1))
	p.DeleteCategory(1)
	fmt.Println(p.GetCategory(1))
	p.RestoreCategory(1)
	fmt.Println(p.GetCategory(1))
	fmt.Println("Server is running...")
	routes.Run()
}
