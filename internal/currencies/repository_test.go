package currencies

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsertCurrency(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{Code: "USD", Symbol: "$", Name: "US Dollar"}
	mock.ExpectQuery("SELECT id FROM currencies").
		WithArgs(currency.Code).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO currencies").
		WithArgs(currency.Code, currency.Symbol, currency.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = CurrencyRepository{db}
	err = p.InsertCurrency(currency)

	if err != nil {
		t.Errorf("error was not expected while inserting currency: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllCurrencies(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{Code: "USD", Symbol: "$", Name: "US Dollar"}
	mock.ExpectQuery("SELECT id, code, symbol, name FROM currencies").
		WillReturnRows(sqlmock.NewRows([]string{"id", "code", "symbol", "name"}).
			AddRow(1, currency.Code, currency.Symbol, currency.Name))

	var p = CurrencyRepository{db}
	r, err := p.GetAllCurrencies()

	if err != nil {
		t.Errorf("error was not expected while getting currencies: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, c := range r {
		if c.Code != currency.Code {
			t.Errorf("expected %s but got %s", currency.Code, c.Code)
		}
	}
}

func TestGetCurrencyByCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{Code: "USD", Symbol: "$", Name: "US Dollar"}
	mock.ExpectQuery("SELECT id, code, symbol, name FROM currencies").
		WithArgs(currency.Code).
		WillReturnRows(sqlmock.NewRows([]string{"id", "code", "symbol", "name"}).
			AddRow(1, currency.Code, currency.Symbol, currency.Name))

	var p = CurrencyRepository{db}
	r, err := p.GetCurrencyByCode(currency.Code)

	if err != nil {
		t.Errorf("error was not expected while getting currency by code: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if r.Code != currency.Code {
		t.Errorf("expected %s but got %s", currency.Code, r.Code)
	}
}

func TestGetCurrencyByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{ID: 1, Code: "USD", Symbol: "$", Name: "US Dollar"}
	mock.ExpectQuery("SELECT id, code, symbol, name FROM currencies").
		WithArgs(currency.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "code", "symbol", "name"}).
			AddRow(currency.ID, currency.Code, currency.Symbol, currency.Name))

	var p = CurrencyRepository{db}
	r, err := p.GetCurrencyByID(currency.ID)

	if err != nil {
		t.Errorf("error was not expected while getting currency by id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if r.ID != currency.ID {
		t.Errorf("expected %d but got %d", currency.ID, r.ID)
	}
}

func TestUpdateCurrency(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{ID: 1, Code: "USD", Symbol: "$", Name: "US Dollar"}
	mock.ExpectExec("UPDATE currencies").
		WithArgs(currency.Code, currency.Symbol, currency.Name, currency.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = CurrencyRepository{db}
	err = p.UpdateCurrency(currency)

	if err != nil {
		t.Errorf("error was not expected while updating currency: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteCurrency(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	currency := Currency{ID: 1, Code: "USD", Symbol: "$", Name: "US Dollar"}
	mock.ExpectExec("DELETE FROM currencies").
		WithArgs(currency.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = CurrencyRepository{db}
	err = p.DeleteCurrency(currency.ID)

	if err != nil {
		t.Errorf("error was not expected while deleting currency: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertExchangeRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse(time.DateOnly, "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	exchangeRate := ExchangeRate{CurrencyID: 1, BaseCurrencyID: 2, Rate: 1.0, Date: twentyTwenty}
	mock.ExpectQuery("SELECT id FROM exchange_rates").
		WithArgs(exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID, exchangeRate.Date).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO exchange_rates").
		WithArgs(exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID, exchangeRate.Rate, exchangeRate.Date).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = CurrencyRepository{db}
	err = p.InsertExchangeRate(exchangeRate)

	if err != nil {
		t.Errorf("error was not expected while inserting exchange rate: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllExchangeRates(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse(time.DateOnly, "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	exchangeRate := ExchangeRate{CurrencyID: 1, BaseCurrencyID: 2, Rate: 1.0, Date: twentyTwenty}
	mock.ExpectQuery("SELECT id, currency_id, base_currency_id, rate, date FROM exchange_rates").
		WillReturnRows(sqlmock.NewRows([]string{"id", "currency_id", "base_currency_id", "rate", "date"}).
			AddRow(1, exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID, exchangeRate.Rate, exchangeRate.Date))

	var p = CurrencyRepository{db}
	r, err := p.GetAllExchangeRates()

	if err != nil {
		t.Errorf("error was not expected while getting exchange rates: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, e := range r {
		if e.CurrencyID != exchangeRate.CurrencyID {
			t.Errorf("expected %d but got %d", exchangeRate.CurrencyID, e.CurrencyID)
		}
	}
}

func TestGetAllExchangeRatesForCurrencies(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse(time.DateOnly, "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	exchangeRate := ExchangeRate{CurrencyID: 1, BaseCurrencyID: 2, Rate: 1.0, Date: twentyTwenty}
	mock.ExpectQuery("SELECT id, currency_id, base_currency_id, rate, date FROM exchange_rates").
		WithArgs(exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "currency_id", "base_currency_id", "rate", "date"}).
			AddRow(1, exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID, exchangeRate.Rate, exchangeRate.Date))

	var p = CurrencyRepository{db}
	r, err := p.GetAllExchangeRatesForCurrencies(exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID)

	if err != nil {
		t.Errorf("error was not expected while getting exchange rates for currencies: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, e := range r {
		if e.CurrencyID != exchangeRate.CurrencyID {
			t.Errorf("expected %d but got %d", exchangeRate.CurrencyID, e.CurrencyID)
		}
	}
}

func TestGetExchangeRateByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse(time.DateOnly, "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	exchangeRate := ExchangeRate{ID: 1, CurrencyID: 1, BaseCurrencyID: 2, Rate: 1.0, Date: twentyTwenty}
	mock.ExpectQuery("SELECT id, currency_id, base_currency_id, rate, date FROM exchange_rates").
		WithArgs(exchangeRate.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "currency_id", "base_currency_id", "rate", "date"}).
			AddRow(exchangeRate.ID, exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID, exchangeRate.Rate, exchangeRate.Date))

	var p = CurrencyRepository{db}
	r, err := p.GetExchangeRateByID(exchangeRate.ID)

	if err != nil {
		t.Errorf("error was not expected while getting exchange rate by id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if r.ID != exchangeRate.ID {
		t.Errorf("expected %d but got %d", exchangeRate.ID, r.ID)
	}
}

func TestGetExchangeRateByCurrenciesAndDate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse(time.DateOnly, "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	exchangeRate := ExchangeRate{CurrencyID: 1, BaseCurrencyID: 2, Rate: 1.0, Date: twentyTwenty}
	mock.ExpectQuery("SELECT id, currency_id, base_currency_id, rate, date FROM exchange_rates").
		WithArgs(exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID, exchangeRate.Date.String()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "currency_id", "base_currency_id", "rate", "date"}).
			AddRow(1, exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID, exchangeRate.Rate, exchangeRate.Date))

	var p = CurrencyRepository{db}
	r, err := p.GetExchangeRateByCurrenciesAndDate(exchangeRate.CurrencyID, exchangeRate.BaseCurrencyID, exchangeRate.Date.String())

	if err != nil {
		t.Errorf("error was not expected while getting exchange rate by currencies and date: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if r.CurrencyID != exchangeRate.CurrencyID {
		t.Errorf("expected %d but got %d", exchangeRate.CurrencyID, r.CurrencyID)
	}
}

func TestUpdateExchangeRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	twentyTwenty, err := time.Parse(time.DateOnly, "2020-01-01")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing time", err)
	}

	exchangeRate := ExchangeRate{ID: 1, CurrencyID: 1, BaseCurrencyID: 2, Rate: 1.0, Date: twentyTwenty}
	mock.ExpectExec("UPDATE exchange_rates").
		WithArgs(exchangeRate.Rate, exchangeRate.Date, exchangeRate.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = CurrencyRepository{db}
	err = p.UpdateExchangeRate(exchangeRate)

	if err != nil {
		t.Errorf("error was not expected while updating exchange rate: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteExchangeRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	exchangeRate := ExchangeRate{ID: 1, CurrencyID: 1, BaseCurrencyID: 2, Rate: 1.0}
	mock.ExpectExec("DELETE FROM exchange_rates").
		WithArgs(exchangeRate.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = CurrencyRepository{db}
	err = p.DeleteExchangeRate(exchangeRate.ID)

	if err != nil {
		t.Errorf("error was not expected while deleting exchange rate: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertCurrencyConversionProvider(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	provider := CurrencyConversionProvider{Name: "Provider", Type: 1, Endpoint: "http://localhost", Enabled: true, Params: "{}", RuntAt: "* * * * *"}
	mock.ExpectQuery("SELECT id FROM currency_conversion_providers").
		WithArgs(provider.Name).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO currency_conversion_providers").
		WithArgs(provider.Name, provider.Type, provider.Endpoint, provider.Enabled, provider.Params, provider.RuntAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = CurrencyRepository{db}
	err = p.InsertCurrencyConversionProvider(provider)

	if err != nil {
		t.Errorf("error was not expected while inserting currency conversion provider: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllCurrencyConversionProviders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	provider := CurrencyConversionProvider{Name: "Provider", Type: 1, Endpoint: "http://localhost", Enabled: true, Params: "{}", RuntAt: "* * * * *"}
	mock.ExpectQuery("SELECT id, name, type, endpoint, enabled, params, runt_at FROM currency_conversion_providers").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "endpoint", "enabled", "params", "runt_at"}).
			AddRow(1, provider.Name, provider.Type, provider.Endpoint, provider.Enabled, provider.Params, provider.RuntAt))

	var p = CurrencyRepository{db}
	r, err := p.GetAllCurrencyConversionProviders()

	if err != nil {
		t.Errorf("error was not expected while getting currency conversion providers: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, c := range r {
		if c.Name != provider.Name {
			t.Errorf("expected %s but got %s", provider.Name, c.Name)
		}
	}
}

func TestGetCurrencyConversionProvidersByCurrencyID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	provider := CurrencyConversionProvider{ID: 1, Name: "Provider", Type: 1, Endpoint: "http://localhost", Enabled: true, Params: "{}", RuntAt: "* * * * *"}
	mock.ExpectQuery("SELECT").
		WithArgs(provider.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "endpoint", "enabled", "params", "runt_at"}).
			AddRow(provider.ID, provider.Name, provider.Type, provider.Endpoint, provider.Enabled, provider.Params, provider.RuntAt))

	var p = CurrencyRepository{db}
	r, err := p.GetCurrencyConversionProvidersByCurrencyID(provider.ID)

	if err != nil {
		t.Errorf("error was not expected while getting currency conversion providers by currency id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, c := range r {
		if c.Name != provider.Name {
			t.Errorf("expected %s but got %s", provider.Name, c.Name)
		}
	}
}

func TestGetCurrencyConversionProviderByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	provider := CurrencyConversionProvider{ID: 1, Name: "Provider", Type: 1, Endpoint: "http://localhost", Enabled: true, Params: "{}", RuntAt: "* * * * *"}
	mock.ExpectQuery("SELECT id, name, type, endpoint, enabled, params, runt_at FROM currency_conversion_providers").
		WithArgs(provider.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "endpoint", "enabled", "params", "runt_at"}).
			AddRow(provider.ID, provider.Name, provider.Type, provider.Endpoint, provider.Enabled, provider.Params, provider.RuntAt))

	var p = CurrencyRepository{db}
	r, err := p.GetCurrencyConversionProviderByID(provider.ID)

	if err != nil {
		t.Errorf("error was not expected while getting currency conversion provider by id: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if r.ID != provider.ID {
		t.Errorf("expected %d but got %d", provider.ID, r.ID)
	}
}

func TestUpdateCurrencyConversionProvider(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	provider := CurrencyConversionProvider{ID: 1, Name: "Provider", Type: 1, Endpoint: "http://localhost", Enabled: true, Params: "{}", RuntAt: "* * * * *"}
	mock.ExpectExec("UPDATE currency_conversion_providers").
		WithArgs(provider.Name, provider.Type, provider.Endpoint, provider.Enabled, provider.Params, provider.RuntAt, provider.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = CurrencyRepository{db}
	err = p.UpdateCurrencyConversionProvider(provider)

	if err != nil {
		t.Errorf("error was not expected while updating currency conversion provider: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteCurrencyConversionProvider(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	provider := CurrencyConversionProvider{ID: 1, Name: "Provider", Type: 1, Endpoint: "http://localhost", Enabled: true, Params: "{}"}
	mock.ExpectExec("DELETE FROM currency_conversion_providers").
		WithArgs(provider.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var p = CurrencyRepository{db}
	err = p.DeleteCurrencyConversionProvider(provider.ID)

	if err != nil {
		t.Errorf("error was not expected while deleting currency conversion provider: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
