package repositories

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/prulloac/fineasy/internal/persistence/entity"
	"github.com/prulloac/fineasy/pkg"
)

type ExchangeRateRepository struct {
	db *sql.DB
}

func NewExchangeRateRepository(db *sql.DB) *ExchangeRateRepository {
	return &ExchangeRateRepository{db}
}

func (e *ExchangeRateRepository) CreateTable() {
	data, _ := os.ReadFile("internal/persistence/schema/exchange_rates.sql")

	if data == nil {
		panic("Error reading accounts schema file!")
	}

	_, err := e.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating exchange rates table!")
		panic(err)
	}
	fmt.Println("Exchange rates table created!")
}

func (e *ExchangeRateRepository) DropTable() {
	_, err := e.db.Exec("DROP TABLE IF EXISTS exchange_rates")
	if err != nil {
		fmt.Println("Error dropping exchange rates table!")
		panic(err)
	}
	fmt.Println("Exchange rates table dropped!")
}

func (e *ExchangeRateRepository) Insert(exchangeRate entity.ExchangeRate) error {
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

func (e *ExchangeRateRepository) GetByCurrencyAndGroupAndTimeFrame(currencyID int, groupID int, timeFrame pkg.Timeframe) ([]entity.ExchangeRate, error) {
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
		currencyID,
		groupID,
		timeFrame.Since,
		timeFrame.Until)
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

func (e *ExchangeRateRepository) GetByCurrencyAndGroupAndDate(currencyID int, groupID int, date time.Time) (entity.ExchangeRate, error) {
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
		currencyID,
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

func (e *ExchangeRateRepository) Update(exchangeRate entity.ExchangeRate) error {
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
