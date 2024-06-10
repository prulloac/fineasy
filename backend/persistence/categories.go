package persistence

import (
	"database/sql"
	"fmt"
	"os"
)

func CreateCategoriesTable(db *sql.DB) {
	data, _ := os.ReadFile("schema/categories.sql")
	_, err := db.Exec(string(data))
	if err != nil {
		panic(err)
	}
	fmt.Println("Categories table created!")
}
