package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/prulloac/fineasy/persistence/entity"
)

type CurrencyRepository struct {
	db *sql.DB
}

func NewCurrencyRepository(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{db}
}

func (c *CurrencyRepository) CreateCurrenciesTable() {
	data, _ := os.ReadFile("persistence/schema/currencies.sql")
	_, err := c.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating currencies table!")
		panic(err)
	}
	fmt.Println("Currencies table created!")
}

func (c *CurrencyRepository) InsertCurrency(currency entity.Currency) error {
	// check if the currency already exists in the database
	var id int
	err := c.db.QueryRow(`
	SELECT
		id 
	FROM currencies 
	WHERE code = $1
	`, currency.Code).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := c.db.Exec(`
		INSERT INTO currencies 
		(code, symbol, name) VALUES ($1, $2, $3)
		`, currency.Code, currency.Symbol, currency.Name)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (c *CurrencyRepository) GetCurrencies() ([]entity.Currency, error) {
	rows, err := c.db.Query(`
	SELECT 
		id, 
		code, 
		symbol, 
		name 
	FROM currencies
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []entity.Currency
	for rows.Next() {
		var currency entity.Currency
		err := rows.Scan(&currency.ID, &currency.Code, &currency.Symbol, &currency.Name)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func (c *CurrencyRepository) GetCurrency(id int) (entity.Currency, error) {
	var currency entity.Currency
	err := c.db.QueryRow(`
	SELECT 
		id, 
		code, 
		symbol, 
		name 
	FROM currencies 
	WHERE id = $1
	`, id).Scan(&currency.ID, &currency.Code, &currency.Symbol, &currency.Name)
	if err != nil {
		return currency, err
	}
	return currency, nil
}

func (c *CurrencyRepository) UpdateCurrency(currency entity.Currency) error {
	_, err := c.db.Exec(`
	UPDATE currencies 
	SET 
		code = $1, 
		symbol = $2, 
		name = $3 
	WHERE id = $4
	`, currency.Code, currency.Symbol, currency.Name, currency.ID)
	if err != nil {
		return err
	}
	return nil
}
