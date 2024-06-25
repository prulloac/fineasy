package currencies

import (
	"encoding/json"
	"fmt"
	"time"
)

type Currency struct {
	ID     int    `json:"id" validate:"required,min=1"`
	Name   string `json:"name" validate:"required,min=2,max=254"`
	Code   string `json:"code" validate:"required,min=2,max=30"`
	Symbol string `json:"symbol" validate:"required,min=1,max=14"`
}

func (c *Currency) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("%+v", c.Name)
	}
	return string(out)
}

type ExchangeRate struct {
	ID             int       `json:"id" validate:"required,min=1"`
	CurrencyID     int       `json:"currency_id" validate:"required,min=1"`
	BaseCurrencyID int       `json:"base_currency_id" validate:"required,min=1"`
	Rate           float64   `json:"rate" validate:"required,min=0.01"`
	Date           time.Time `json:"date" validate:"required,past_time"`
}

func (e *ExchangeRate) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("%+v", e.ID)
	}
	return string(out)
}

type CurrencyConversionProvider struct {
	ID       int    `json:"id" validate:"required,min=1"`
	Name     string `json:"name" validate:"required,min=2,max=254"`
	Type     int    `json:"type" validate:"required"`
	Endpoint string `json:"endpoint" validate:"required,url"`
	Enabled  bool   `json:"enabled" validate:"required,boolean"`
	Params   string `json:"params" validate:"required,json"`
	RuntAt   string `json:"runt_at" validate:"required,cron"`
}

func (e *CurrencyConversionProvider) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("%+v", e.ID)
	}
	return string(out)
}
