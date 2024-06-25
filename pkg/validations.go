package pkg

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func PastTime(fl validator.FieldLevel) bool {
	date := fl.Field().Interface().(time.Time)
	return date.Before(time.Now())
}

func ValidateStruct(s interface{}) error {
	validate := validator.New()
	validate.RegisterValidation("past_time", PastTime)
	return validate.Struct(s)
}
