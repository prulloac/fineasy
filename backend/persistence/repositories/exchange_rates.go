package repositories

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/prulloac/fineasy/persistence/entity"
)

type ExchangeRateRepository struct {
	db *sql.DB
}

func NewExchangeRateRepository(db *sql.DB) *ExchangeRateRepository {
	return &ExchangeRateRepository{db}
}

func (e *ExchangeRateRepository) CreateExchangeRatesTable() {
	data, _ := os.ReadFile("persistence/schema/exchange_rates.sql")
	_, err := e.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating exchange rates table!")
		panic(err)
	}
	fmt.Println("Exchange rates table created!")
}

func (e *ExchangeRateRepository) InsertExchangeRate(exchangeRate entity.ExchangeRate) error {
	// check if the exchange rate already exists
	var id int
	err := e.db.QueryRow(`
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
		_, err := e.db.Exec(`
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

func (e *ExchangeRateRepository) GetExchangeRates(currency entity.Currency, groupID int, since time.Time, until time.Time) ([]entity.ExchangeRate, error) {
	rows, err := e.db.Query(`
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
		groupID,
		since,
		until)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exchangeRates []entity.ExchangeRate
	for rows.Next() {
		var exchangeRate entity.ExchangeRate
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

func (e *ExchangeRateRepository) GetExchangeRate(currency entity.Currency, groupID int, date time.Time) (entity.ExchangeRate, error) {
	var exchangeRate entity.ExchangeRate
	err := e.db.QueryRow(`
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
		groupID,
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

func (e *ExchangeRateRepository) UpdateExchangeRate(exchangeRate entity.ExchangeRate) error {
	_, err := e.db.Exec(`
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
