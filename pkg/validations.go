package pkg

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func PastTime(fl validator.FieldLevel) bool {
	date := fl.Field().Interface().(time.Time)
	return date.Before(time.Now())
}
