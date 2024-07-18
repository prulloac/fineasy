package validations

import (
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func PastTime(fl validator.FieldLevel) bool {
	if fl.Field().Type().String() == "time.Time" {
		date := fl.Field().Interface().(time.Time)
		return date.Before(time.Now())
	}
	if fl.Field().Type().String() == "sql.NullTime" {
		date := fl.Field().Interface().(sql.NullTime)
		if date.Valid {
			return date.Time.Before(time.Now())
		}
		return false
	}
	return false
}

func UUID7(fl validator.FieldLevel) bool {
	v, err := uuid.Parse(fl.Field().String())
	if err != nil {
		return false
	}
	if v.Version() == 7 {
		return true
	}
	return false
}

func Date(fl validator.FieldLevel) bool {
	if fl.Field().Type().String() == "time.Time" {
		date := fl.Field().Interface().(time.Time)
		return date.Hour() == 0 && date.Minute() == 0 && date.Second() == 0
	}
	if fl.Field().Type().String() == "sql.NullTime" {
		date := fl.Field().Interface().(sql.NullTime)
		if date.Valid {
			return date.Time.Hour() == 0 && date.Time.Minute() == 0 && date.Time.Second() == 0
		}
		return false
	}
	if fl.Field().Type().String() == "string" {
		date := fl.Field().Interface().(string)
		t, err := time.Parse(time.DateOnly, date)
		if err != nil {
			return false
		}
		return t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0
	}
	return false
}
