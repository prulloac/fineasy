package pkg

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func PastTime(fl validator.FieldLevel) bool {
	date := fl.Field().Interface().(time.Time)
	return date.Before(time.Now())
}

func UUID7(fl validator.FieldLevel) bool {
	v, err := uuid.Parse(fl.Field().String())
	if err != nil {
		return false
	}
	if v.Version() != 7 {
		return false
	}
	return true
}

func ValidateStruct(s interface{}) error {
	validate := validator.New()
	validate.RegisterValidation("past_time", PastTime)
	validate.RegisterValidation("uuid7", UUID7)
	return validate.Struct(s)
}
