package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Currency struct {
	ID     int
	Code   string
	Symbol string
	Name   string
}

type ExchangeRate struct {
	CurrencyID int
	GroupID    int
	Rate       float64
	Date       time.Time
}

func (p *Persistence) CreateCurrenciesTable() {
	data, _ := os.ReadFile("persistence/schema/currencies.sql")
	_, err := p.db.Exec(string(data))
	if err != nil {
		panic(err)
	}
	fmt.Println("Currencies table created!")
}

func (p *Persistence) InsertCurrency(currency Currency) error {
	// check if the currency already exists in the database
	var id int
	err := p.db.QueryRow(`
	SELECT
		id 
	FROM currencies 
	WHERE code = $1
	`, currency.Code).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := p.db.Exec(`
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

func (p *Persistence) GetCurrencies() ([]Currency, error) {
	rows, err := p.db.Query(`
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

	var currencies []Currency
	for rows.Next() {
		var currency Currency
		err := rows.Scan(&currency.ID, &currency.Code, &currency.Symbol, &currency.Name)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func (p *Persistence) GetCurrency(id int) (Currency, error) {
	var currency Currency
	err := p.db.QueryRow(`
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

func (p *Persistence) UpdateCurrency(currency Currency) error {
	_, err := p.db.Exec(`
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

func (p *Persistence) InsertExchangeRate(exchangeRate ExchangeRate) error {
	_, err := p.db.Exec(`
	INSERT INTO default_exchange_rates 
	(currency_id, group_id, rate, date) VALUES ($1, $2, $3, $4)
	`, exchangeRate.CurrencyID, exchangeRate.GroupID, exchangeRate.Rate, exchangeRate.Date)
	if err != nil {
		return err
	}
	return nil
}

func (p *Persistence) GetExchangeRates(currency Currency, group_id int, since time.Time, until time.Time) ([]ExchangeRate, error) {
	rows, err := p.db.Query(`
	SELECT 
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
	`, currency.ID, group_id, since, until)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exchangeRates []ExchangeRate
	for rows.Next() {
		var exchangeRate ExchangeRate
		err := rows.Scan(&exchangeRate.CurrencyID, &exchangeRate.Rate, &exchangeRate.Date, &exchangeRate.GroupID)
		if err != nil {
			return nil, err
		}
		exchangeRates = append(exchangeRates, exchangeRate)
	}
	return exchangeRates, nil
}

func (p *Persistence) GetExchangeRate(currency Currency, group_id int, date time.Time) (ExchangeRate, error) {
	var exchangeRate ExchangeRate
	err := p.db.QueryRow(`
	SELECT 
		currency_id, 
		group_id,
		rate, 
		date 
	FROM default_exchange_rates 
	WHERE 1=1
	AND currency_id = $1 
	AND group_id = $2
	AND date = $3
	`, currency.ID, group_id, date).
		Scan(&exchangeRate.CurrencyID, &exchangeRate.GroupID, &exchangeRate.Rate, &exchangeRate.Date)
	if err != nil {
		return exchangeRate, err
	}
	return exchangeRate, nil
}

func (p *Persistence) UpdateExchangeRate(exchangeRate ExchangeRate) error {
	_, err := p.db.Exec(`
	UPDATE exchange_rates 
	SET rate = $1 
	WHERE 1=1
	AND currency_id = $2 
	AND group_id = $3
	AND date = $4
	`, exchangeRate.Rate, exchangeRate.CurrencyID, exchangeRate.GroupID, exchangeRate.Date)
	if err != nil {
		return err
	}
	return nil
}
