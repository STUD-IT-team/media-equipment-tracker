package validate

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func setup() {
	once.Do(func() {
		validate = validator.New()
	})
}

func ValidateStruct(i interface{}) error {
	setup()
	return validate.Struct(i)
}
