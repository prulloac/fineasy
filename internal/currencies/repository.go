package currencies

import (
	"database/sql"

	"github.com/prulloac/fineasy/pkg"
)

type CurrencyRepository struct {
	db *sql.DB
}

func NewCurrencyRepository(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{db}
}

func (c *CurrencyRepository) InsertCurrency(currency Currency) error {
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

func (c *CurrencyRepository) GetAllCurrencies() ([]Currency, error) {
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

	var currencies []Currency
	for rows.Next() {
		var currency Currency
		err := rows.Scan(&currency.ID, &currency.Code, &currency.Symbol, &currency.Name)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(&currency)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func (c *CurrencyRepository) GetCurrencyByCode(code string) (Currency, error) {
	var currency Currency
	err := c.db.QueryRow(`
	SELECT 
		id, 
		code, 
		symbol, 
		name 
	FROM currencies
	WHERE code = $1
	`, code).Scan(&currency.ID, &currency.Code, &currency.Symbol, &currency.Name)
	if err != nil {
		return Currency{}, err
	}
	err = pkg.ValidateStruct(&currency)
	if err != nil {
		return Currency{}, err
	}
	return currency, nil
}

func (c *CurrencyRepository) GetCurrencyByID(id int) (Currency, error) {
	var currency Currency
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
		return Currency{}, err
	}
	err = pkg.ValidateStruct(&currency)
	if err != nil {
		return Currency{}, err
	}
	return currency, nil
}

func (c *CurrencyRepository) UpdateCurrency(currency Currency) error {
	err := pkg.ValidateStruct(&currency)
	if err != nil {
		return err
	}
	_, err = c.db.Exec(`
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

func (c *CurrencyRepository) DeleteCurrency(id int) error {
	_, err := c.db.Exec(`
	DELETE FROM currencies 
	WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CurrencyRepository) InsertExchangeRate(exchangeRate ExchangeRate) error {
	// check if the exchange rate already exists in the database
	var id int
	err := c.db.QueryRow(`
	SELECT
		id
	FROM exchange_rates
	WHERE currency_id = $1 AND base_currency_id = $2 AND date = $3
	`, exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID, exchangeRate.Date).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := c.db.Exec(`
		INSERT INTO exchange_rates 
		(currency_id, base_currency_id, rate, date) VALUES ($1, $2, $3, $4)
		`, exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID, exchangeRate.Rate, exchangeRate.Date)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CurrencyRepository) GetAllExchangeRates() ([]ExchangeRate, error) {
	rows, err := c.db.Query(`
	SELECT 
		id, 
		currency_id, 
		base_currency_id, 
		rate, 
		date 
	FROM exchange_rates
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exchangeRates []ExchangeRate
	for rows.Next() {
		var exchangeRate ExchangeRate
		err := rows.Scan(&exchangeRate.ID, &exchangeRate.CurrencyID, &exchangeRate.BaseCurrencyID, &exchangeRate.Rate, &exchangeRate.Date)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(&exchangeRate)
		if err != nil {
			return nil, err
		}
		exchangeRates = append(exchangeRates, exchangeRate)
	}
	return exchangeRates, nil
}

func (c *CurrencyRepository) GetAllExchangeRatesForCurrencies(currencyID int, baseCurrencyID int) ([]ExchangeRate, error) {
	rows, err := c.db.Query(`
	SELECT 
		id, 
		currency_id, 
		base_currency_id, 
		rate, 
		date 
	FROM exchange_rates
	WHERE currency_id = $1 AND base_currency_id = $2
	`, currencyID, baseCurrencyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exchangeRates []ExchangeRate
	for rows.Next() {
		var exchangeRate ExchangeRate
		err := rows.Scan(&exchangeRate.ID, &exchangeRate.CurrencyID, &exchangeRate.BaseCurrencyID, &exchangeRate.Rate, &exchangeRate.Date)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(&exchangeRate)
		if err != nil {
			return nil, err
		}
		exchangeRates = append(exchangeRates, exchangeRate)
	}
	return exchangeRates, nil
}

func (c *CurrencyRepository) GetExchangeRateByID(id int) (ExchangeRate, error) {
	var exchangeRate ExchangeRate
	err := c.db.QueryRow(`
	SELECT 
		id, 
		currency_id, 
		base_currency_id, 
		rate, 
		date 
	FROM exchange_rates
	WHERE id = $1
	`, id).Scan(&exchangeRate.ID, &exchangeRate.CurrencyID, &exchangeRate.BaseCurrencyID, &exchangeRate.Rate, &exchangeRate.Date)
	if err != nil {
		return ExchangeRate{}, err
	}
	err = pkg.ValidateStruct(&exchangeRate)
	if err != nil {
		return ExchangeRate{}, err
	}
	return exchangeRate, nil
}

func (c *CurrencyRepository) GetExchangeRateByCurrenciesAndDate(currencyID int, baseCurrencyID int, date string) (ExchangeRate, error) {
	var exchangeRate ExchangeRate
	err := c.db.QueryRow(`
	SELECT 
		id, 
		currency_id, 
		base_currency_id, 
		rate, 
		date 
	FROM exchange_rates
	WHERE currency_id = $1 AND base_currency_id = $2 AND date = $3
	`, currencyID, baseCurrencyID, date).Scan(&exchangeRate.ID, &exchangeRate.CurrencyID, &exchangeRate.BaseCurrencyID, &exchangeRate.Rate, &exchangeRate.Date)
	if err != nil {
		return ExchangeRate{}, err
	}
	err = pkg.ValidateStruct(&exchangeRate)
	if err != nil {
		return ExchangeRate{}, err
	}
	return exchangeRate, nil
}

func (c *CurrencyRepository) UpdateExchangeRate(exchangeRate ExchangeRate) error {
	err := pkg.ValidateStruct(&exchangeRate)
	if err != nil {
		return err
	}
	_, err = c.db.Exec(`
	UPDATE exchange_rates 
	SET 
		rate = $1, 
		date = $2 
	WHERE id = $3
	`, exchangeRate.Rate, exchangeRate.Date, exchangeRate.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CurrencyRepository) DeleteExchangeRate(id int) error {
	_, err := c.db.Exec(`
	DELETE FROM exchange_rates 
	WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CurrencyRepository) InsertCurrencyConversionProvider(provider CurrencyConversionProvider) error {
	// check if the provider already exists in the database
	var id int
	err := c.db.QueryRow(`
	SELECT
		id
	FROM currency_conversion_providers
	WHERE name = $1
	`, provider.Name).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := c.db.Exec(`
		INSERT INTO currency_conversion_providers 
		(name, type, endpoint, enabled, params, runt_at) VALUES ($1, $2, $3, $4, $5, $6)
		`, provider.Name, provider.Type, provider.Endpoint, provider.Enabled, provider.Params, provider.RuntAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CurrencyRepository) GetAllCurrencyConversionProviders() ([]CurrencyConversionProvider, error) {
	rows, err := c.db.Query(`
	SELECT 
		id, 
		name, 
		type, 
		endpoint, 
		enabled, 
		params, 
		runt_at 
	FROM currency_conversion_providers
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []CurrencyConversionProvider
	for rows.Next() {
		var provider CurrencyConversionProvider
		err := rows.Scan(&provider.ID, &provider.Name, &provider.Type, &provider.Endpoint, &provider.Enabled, &provider.Params, &provider.RuntAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(&provider)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

func (c *CurrencyRepository) GetCurrencyConversionProvidersByCurrencyID(currencyID int) ([]CurrencyConversionProvider, error) {
	rows, err := c.db.Query(`
	SELECT 
		ccp.id, 
		ccp.name, 
		ccp.type, 
		ccp.endpoint, 
		ccp.enabled, 
		ccp.params, 
		ccp.runt_at 
	FROM currency_conversion_providers ccp
	JOIN currency_conversion_provider_currencies ccpc ON ccp.id = ccpc.currency_conversion_provider_id
	WHERE ccpc.currency_id = $1
	`, currencyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []CurrencyConversionProvider
	for rows.Next() {
		var provider CurrencyConversionProvider
		err := rows.Scan(&provider.ID, &provider.Name, &provider.Type, &provider.Endpoint, &provider.Enabled, &provider.Params, &provider.RuntAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(&provider)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

func (c *CurrencyRepository) GetCurrencyConversionProviderByID(id int) (CurrencyConversionProvider, error) {
	var provider CurrencyConversionProvider
	err := c.db.QueryRow(`
	SELECT 
		id, 
		name, 
		type, 
		endpoint, 
		enabled, 
		params, 
		runt_at 
	FROM currency_conversion_providers
	WHERE id = $1
	`, id).Scan(&provider.ID, &provider.Name, &provider.Type, &provider.Endpoint, &provider.Enabled, &provider.Params, &provider.RuntAt)
	if err != nil {
		return CurrencyConversionProvider{}, err
	}
	err = pkg.ValidateStruct(&provider)
	if err != nil {
		return CurrencyConversionProvider{}, err
	}
	return provider, nil
}

func (c *CurrencyRepository) UpdateCurrencyConversionProvider(provider CurrencyConversionProvider) error {
	err := pkg.ValidateStruct(&provider)
	if err != nil {
		return err
	}
	_, err = c.db.Exec(`
	UPDATE currency_conversion_providers 
	SET 
		name = $1, 
		type = $2, 
		endpoint = $3, 
		enabled = $4, 
		params = $5, 
		runt_at = $6 
	WHERE id = $7
	`, provider.Name, provider.Type, provider.Endpoint, provider.Enabled, provider.Params, provider.RuntAt, provider.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CurrencyRepository) DeleteCurrencyConversionProvider(id int) error {
	_, err := c.db.Exec(`
	DELETE FROM currency_conversion_providers 
	WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}
