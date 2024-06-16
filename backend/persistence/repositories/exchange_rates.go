package repositories

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	. "github.com/prulloac/fineasy/persistence/entity"
)

type ExchangeRateRepository struct {
	DB *sql.DB
}

func (e *ExchangeRateRepository) CreateExchangeRatesTable() {
	data, _ := os.ReadFile("persistence/schema/exchange_rates.sql")
	_, err := e.DB.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating exchange rates table!")
		panic(err)
	}
	fmt.Println("Exchange rates table created!")
}

func (e *ExchangeRateRepository) InsertExchangeRate(exchangeRate ExchangeRate) error {
	// check if the exchange rate already exists
	var id int
	err := e.DB.QueryRow(`
	SELECT
		id
	FROM exchange_rates
	WHERE currency_id = $1
	AND group_id = $2
	AND date = $3
	`,
		exchangeRate.CurrencyID,
		exchangeRate.GroupID,
		exchangeRate.Date).
		Scan(&id)

	if err == sql.ErrNoRows {
		_, err := e.DB.Exec(`
		INSERT INTO exchange_rates 
		(currency_id, group_id, rate, date) VALUES ($1, $2, $3, $4)
		`,
			exchangeRate.CurrencyID,
			exchangeRate.GroupID,
			exchangeRate.Rate,
			exchangeRate.Date)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func (e *ExchangeRateRepository) GetExchangeRates(currency Currency, group_id int, since time.Time, until time.Time) ([]ExchangeRate, error) {
	rows, err := e.DB.Query(`
	SELECT 
		id,
		currency_id, 
		rate, 
		date,
		group_id
	FROM default_exchange_rates 
	WHERE 1=1
	AND currency_id = $1
	AND group_id = $2 
	AND date >= $3 
	AND date <= $4
	`,
		currency.ID,
		group_id,
		since,
		until)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exchangeRates []ExchangeRate
	for rows.Next() {
		var exchangeRate ExchangeRate
		err := rows.Scan(
			&exchangeRate.ID,
			&exchangeRate.CurrencyID,
			&exchangeRate.Rate,
			&exchangeRate.Date,
			&exchangeRate.GroupID)
		if err != nil {
			return nil, err
		}
		exchangeRates = append(exchangeRates, exchangeRate)
	}
	return exchangeRates, nil
}

func (e *ExchangeRateRepository) GetExchangeRate(currency Currency, group_id int, date time.Time) (ExchangeRate, error) {
	var exchangeRate ExchangeRate
	err := e.DB.QueryRow(`
	SELECT 
		currency_id, 
		group_id,
		rate, 
		date 
	FROM exchange_rates 
	WHERE 1=1
	AND currency_id = $1 
	AND group_id = $2
	AND date = $3
	`,
		currency.ID,
		group_id,
		date).
		Scan(
			&exchangeRate.CurrencyID,
			&exchangeRate.GroupID,
			&exchangeRate.Rate,
			&exchangeRate.Date)
	if err != nil {
		return exchangeRate, err
	}
	return exchangeRate, nil
}

func (e *ExchangeRateRepository) UpdateExchangeRate(exchangeRate ExchangeRate) error {
	_, err := e.DB.Exec(`
	UPDATE exchange_rates 
	SET rate = $1 
	WHERE id = $2
	`,
		exchangeRate.Rate,
		exchangeRate.ID)
	if err != nil {
		return err
	}
	return nil
}
