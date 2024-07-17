package pkg

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
		if !date.Valid {
			return date.Time.Before(time.Now())
		}
		return true
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
